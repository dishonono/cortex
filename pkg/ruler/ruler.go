package ruler

import (
	"context"
	"flag"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/rulefmt"
	"github.com/prometheus/prometheus/notifier"
	promRules "github.com/prometheus/prometheus/rules"
	"github.com/prometheus/prometheus/util/strutil"
	"github.com/weaveworks/common/user"
	"golang.org/x/sync/errgroup"

	"github.com/cortexproject/cortex/pkg/cortexpb"
	"github.com/cortexproject/cortex/pkg/ring"
	"github.com/cortexproject/cortex/pkg/ring/kv"
	"github.com/cortexproject/cortex/pkg/ruler/rulespb"
	"github.com/cortexproject/cortex/pkg/ruler/rulestore"
	"github.com/cortexproject/cortex/pkg/tenant"
	"github.com/cortexproject/cortex/pkg/util"
	"github.com/cortexproject/cortex/pkg/util/concurrency"
	"github.com/cortexproject/cortex/pkg/util/flagext"
	"github.com/cortexproject/cortex/pkg/util/grpcclient"
	util_log "github.com/cortexproject/cortex/pkg/util/log"
	"github.com/cortexproject/cortex/pkg/util/services"
	"github.com/cortexproject/cortex/pkg/util/validation"
)

var (
	supportedShardingStrategies = []string{util.ShardingStrategyDefault, util.ShardingStrategyShuffle}

	// Validation errors.
	errInvalidShardingStrategy = errors.New("invalid sharding strategy")
	errInvalidTenantShardSize  = errors.New("invalid tenant shard size, the value must be greater than 0")
)

const (
	// ringKey is the key under which we store the rulers ring in the KVStore.
	ringKey = "ring"

	// Number of concurrent group list and group loads operations.
	loadRulesConcurrency  = 10
	fetchRulesConcurrency = 16

	rulerSyncReasonInitial    = "initial"
	rulerSyncReasonPeriodic   = "periodic"
	rulerSyncReasonRingChange = "ring-change"

	// Limit errors
	errMaxRuleGroupsPerUserLimitExceeded        = "per-user rule groups limit (limit: %d actual: %d) exceeded"
	errMaxRulesPerRuleGroupPerUserLimitExceeded = "per-user rules per rule group limit (limit: %d actual: %d) exceeded"

	// errors
	errListAllUser = "unable to list the ruler users"

	alertingRuleFilter  string = "alert"
	recordingRuleFilter string = "record"
)

type DisabledRuleGroupErr struct {
	Message string
}

func (e *DisabledRuleGroupErr) Error() string {
	return e.Message
}

// Config is the configuration for the recording rules server.
type Config struct {
	// This is used for template expansion in alerts; must be a valid URL.
	ExternalURL flagext.URLValue `yaml:"external_url"`
	// Labels to add to all alerts
	ExternalLabels labels.Labels `yaml:"external_labels,omitempty" doc:"nocli|description=Labels to add to all alerts."`
	// GRPC Client configuration.
	ClientTLSConfig grpcclient.Config `yaml:"ruler_client"`
	// How frequently to evaluate rules by default.
	EvaluationInterval time.Duration `yaml:"evaluation_interval"`
	// How frequently to poll for updated rules.
	PollInterval time.Duration `yaml:"poll_interval"`
	// Path to store rule files for prom manager.
	RulePath string `yaml:"rule_path"`

	// URL of the Alertmanager to send notifications to.
	// If your are configuring the ruler to send to a Cortex Alertmanager,
	// ensure this includes any path set in the Alertmanager external URL.
	AlertmanagerURL string `yaml:"alertmanager_url"`
	// Whether to use DNS SRV records to discover Alertmanager.
	AlertmanagerDiscovery bool `yaml:"enable_alertmanager_discovery"`
	// How long to wait between refreshing the list of Alertmanager based on DNS service discovery.
	AlertmanagerRefreshInterval time.Duration `yaml:"alertmanager_refresh_interval"`
	// Enables the ruler notifier to use the Alertmananger V2 API.
	AlertmanangerEnableV2API bool `yaml:"enable_alertmanager_v2"`
	// Capacity of the queue for notifications to be sent to the Alertmanager.
	NotificationQueueCapacity int `yaml:"notification_queue_capacity"`
	// HTTP timeout duration when sending notifications to the Alertmanager.
	NotificationTimeout time.Duration `yaml:"notification_timeout"`
	// Client configs for interacting with the Alertmanager
	Notifier NotifierConfig `yaml:"alertmanager_client"`

	// Max time to tolerate outage for restoring "for" state of alert.
	OutageTolerance time.Duration `yaml:"for_outage_tolerance"`
	// Minimum duration between alert and restored "for" state. This is maintained only for alerts with configured "for" time greater than grace period.
	ForGracePeriod time.Duration `yaml:"for_grace_period"`
	// Minimum amount of time to wait before resending an alert to Alertmanager.
	ResendDelay time.Duration `yaml:"resend_delay"`

	// Enable sharding rule groups.
	EnableSharding   bool          `yaml:"enable_sharding"`
	ShardingStrategy string        `yaml:"sharding_strategy"`
	SearchPendingFor time.Duration `yaml:"search_pending_for"`
	Ring             RingConfig    `yaml:"ring"`
	FlushCheckPeriod time.Duration `yaml:"flush_period"`

	EnableAPI bool `yaml:"enable_api"`

	EnabledTenants  flagext.StringSliceCSV `yaml:"enabled_tenants"`
	DisabledTenants flagext.StringSliceCSV `yaml:"disabled_tenants"`

	RingCheckPeriod time.Duration `yaml:"-"`

	EnableQueryStats      bool `yaml:"query_stats_enabled"`
	DisableRuleGroupLabel bool `yaml:"disable_rule_group_label"`
}

// Validate config and returns error on failure
func (cfg *Config) Validate(limits validation.Limits, log log.Logger) error {
	if !util.StringsContain(supportedShardingStrategies, cfg.ShardingStrategy) {
		return errInvalidShardingStrategy
	}

	if cfg.ShardingStrategy == util.ShardingStrategyShuffle && limits.RulerTenantShardSize <= 0 {
		return errInvalidTenantShardSize
	}

	if err := cfg.ClientTLSConfig.Validate(log); err != nil {
		return errors.Wrap(err, "invalid ruler gRPC client config")
	}
	return nil
}

// RegisterFlags adds the flags required to config this to the given FlagSet
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.ClientTLSConfig.RegisterFlagsWithPrefix("ruler.client", f)
	cfg.Ring.RegisterFlags(f)
	cfg.Notifier.RegisterFlags(f)

	// Deprecated Flags that will be maintained to avoid user disruption

	//lint:ignore faillint Need to pass the global logger like this for warning on deprecated methods
	flagext.DeprecatedFlag(f, "ruler.client-timeout", "This flag has been renamed to ruler.configs.client-timeout", util_log.Logger)
	//lint:ignore faillint Need to pass the global logger like this for warning on deprecated methods
	flagext.DeprecatedFlag(f, "ruler.group-timeout", "This flag is no longer functional.", util_log.Logger)
	//lint:ignore faillint Need to pass the global logger like this for warning on deprecated methods
	flagext.DeprecatedFlag(f, "ruler.num-workers", "This flag is no longer functional. For increased concurrency horizontal sharding is recommended", util_log.Logger)

	cfg.ExternalURL.URL, _ = url.Parse("") // Must be non-nil
	f.Var(&cfg.ExternalURL, "ruler.external.url", "URL of alerts return path.")
	f.DurationVar(&cfg.EvaluationInterval, "ruler.evaluation-interval", 1*time.Minute, "How frequently to evaluate rules")
	f.DurationVar(&cfg.PollInterval, "ruler.poll-interval", 1*time.Minute, "How frequently to poll for rule changes")

	f.StringVar(&cfg.AlertmanagerURL, "ruler.alertmanager-url", "", "Comma-separated list of URL(s) of the Alertmanager(s) to send notifications to. Each Alertmanager URL is treated as a separate group in the configuration. Multiple Alertmanagers in HA per group can be supported by using DNS resolution via -ruler.alertmanager-discovery.")
	f.BoolVar(&cfg.AlertmanagerDiscovery, "ruler.alertmanager-discovery", false, "Use DNS SRV records to discover Alertmanager hosts.")
	f.DurationVar(&cfg.AlertmanagerRefreshInterval, "ruler.alertmanager-refresh-interval", 1*time.Minute, "How long to wait between refreshing DNS resolutions of Alertmanager hosts.")
	f.BoolVar(&cfg.AlertmanangerEnableV2API, "ruler.alertmanager-use-v2", false, "If enabled requests to Alertmanager will utilize the V2 API.")
	f.IntVar(&cfg.NotificationQueueCapacity, "ruler.notification-queue-capacity", 10000, "Capacity of the queue for notifications to be sent to the Alertmanager.")
	f.DurationVar(&cfg.NotificationTimeout, "ruler.notification-timeout", 10*time.Second, "HTTP timeout duration when sending notifications to the Alertmanager.")

	f.DurationVar(&cfg.SearchPendingFor, "ruler.search-pending-for", 5*time.Minute, "Time to spend searching for a pending ruler when shutting down.")
	f.BoolVar(&cfg.EnableSharding, "ruler.enable-sharding", false, "Distribute rule evaluation using ring backend")
	f.StringVar(&cfg.ShardingStrategy, "ruler.sharding-strategy", util.ShardingStrategyDefault, fmt.Sprintf("The sharding strategy to use. Supported values are: %s.", strings.Join(supportedShardingStrategies, ", ")))
	f.DurationVar(&cfg.FlushCheckPeriod, "ruler.flush-period", 1*time.Minute, "Period with which to attempt to flush rule groups.")
	f.StringVar(&cfg.RulePath, "ruler.rule-path", "/rules", "file path to store temporary rule files for the prometheus rule managers")
	f.BoolVar(&cfg.EnableAPI, "experimental.ruler.enable-api", false, "Enable the ruler api")
	f.DurationVar(&cfg.OutageTolerance, "ruler.for-outage-tolerance", time.Hour, `Max time to tolerate outage for restoring "for" state of alert.`)
	f.DurationVar(&cfg.ForGracePeriod, "ruler.for-grace-period", 10*time.Minute, `Minimum duration between alert and restored "for" state. This is maintained only for alerts with configured "for" time greater than grace period.`)
	f.DurationVar(&cfg.ResendDelay, "ruler.resend-delay", time.Minute, `Minimum amount of time to wait before resending an alert to Alertmanager.`)

	f.Var(&cfg.EnabledTenants, "ruler.enabled-tenants", "Comma separated list of tenants whose rules this ruler can evaluate. If specified, only these tenants will be handled by ruler, otherwise this ruler can process rules from all tenants. Subject to sharding.")
	f.Var(&cfg.DisabledTenants, "ruler.disabled-tenants", "Comma separated list of tenants whose rules this ruler cannot evaluate. If specified, a ruler that would normally pick the specified tenant(s) for processing will ignore them instead. Subject to sharding.")

	f.BoolVar(&cfg.EnableQueryStats, "ruler.query-stats-enabled", false, "Report the wall time for ruler queries to complete as a per user metric and as an info level log message.")
	f.BoolVar(&cfg.DisableRuleGroupLabel, "ruler.disable-rule-group-label", false, "Disable the rule_group label on exported metrics")

	cfg.RingCheckPeriod = 5 * time.Second
}

// MultiTenantManager is the interface of interaction with a Manager that is tenant aware.
type MultiTenantManager interface {
	// SyncRuleGroups is used to sync the Manager with rules from the RuleStore.
	// If existing user is missing in the ruleGroups map, its ruler manager will be stopped.
	SyncRuleGroups(ctx context.Context, ruleGroups map[string]rulespb.RuleGroupList)
	// GetRules fetches rules for a particular tenant (userID).
	GetRules(userID string) []*promRules.Group
	// Stop stops all Manager components.
	Stop()
	// ValidateRuleGroup validates a rulegroup
	ValidateRuleGroup(rulefmt.RuleGroup) []error
}

// Ruler evaluates rules.
//
//	+---------------------------------------------------------------+
//	|                                                               |
//	|                   Query       +-------------+                 |
//	|            +------------------>             |                 |
//	|            |                  |    Store    |                 |
//	|            | +----------------+             |                 |
//	|            | |     Rules      +-------------+                 |
//	|            | |                                                |
//	|            | |                                                |
//	|            | |                                                |
//	|       +----+-v----+   Filter  +------------+                  |
//	|       |           +----------->            |                  |
//	|       |   Ruler   |           |    Ring    |                  |
//	|       |           <-----------+            |                  |
//	|       +-------+---+   Rules   +------------+                  |
//	|               |                                               |
//	|               |                                               |
//	|               |                                               |
//	|               |    Load      +-----------------+              |
//	|               +-------------->                 |              |
//	|                              |     Manager     |              |
//	|                              |                 |              |
//	|                              +-----------------+              |
//	|                                                               |
//	+---------------------------------------------------------------+
type Ruler struct {
	services.Service

	cfg        Config
	lifecycler *ring.BasicLifecycler
	ring       *ring.Ring
	store      rulestore.RuleStore
	manager    MultiTenantManager
	limits     RulesLimits

	subservices        *services.Manager
	subservicesWatcher *services.FailureWatcher

	// Pool of clients used to connect to other ruler replicas.
	clientsPool ClientsPool

	ringCheckErrors            prometheus.Counter
	rulerSync                  *prometheus.CounterVec
	ruleGroupStoreLoadDuration prometheus.Gauge
	ruleGroupSyncDuration      prometheus.Gauge

	allowedTenants *util.AllowedTenants

	registry prometheus.Registerer
	logger   log.Logger
}

// NewRuler creates a new ruler from a distributor and chunk store.
func NewRuler(cfg Config, manager MultiTenantManager, reg prometheus.Registerer, logger log.Logger, ruleStore rulestore.RuleStore, limits RulesLimits) (*Ruler, error) {
	return newRuler(cfg, manager, reg, logger, ruleStore, limits, newRulerClientPool(cfg.ClientTLSConfig, logger, reg))
}

func (r *Ruler) GetRing() *ring.Ring {
	return r.ring
}

func (r *Ruler) GetLimits() RulesLimits {
	return r.limits
}

func (r *Ruler) GetLifecycler() *ring.BasicLifecycler {
	return r.lifecycler
}
func (r *Ruler) GetClientsPool() ClientsPool {
	return r.clientsPool
}

func newRuler(cfg Config, manager MultiTenantManager, reg prometheus.Registerer, logger log.Logger, ruleStore rulestore.RuleStore, limits RulesLimits, clientPool ClientsPool) (*Ruler, error) {
	ruler := &Ruler{
		cfg:            cfg,
		store:          ruleStore,
		manager:        manager,
		registry:       reg,
		logger:         logger,
		limits:         limits,
		clientsPool:    clientPool,
		allowedTenants: util.NewAllowedTenants(cfg.EnabledTenants, cfg.DisabledTenants),

		ringCheckErrors: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "cortex_ruler_ring_check_errors_total",
			Help: "Number of errors that have occurred when checking the ring for ownership",
		}),

		rulerSync: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "cortex_ruler_sync_rules_total",
			Help: "Total number of times the ruler sync operation triggered.",
		}, []string{"reason"}),

		ruleGroupStoreLoadDuration: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Name: "cortex_ruler_rule_group_load_duration_seconds",
			Help: "Time taken to load rule groups from storage",
		}),

		ruleGroupSyncDuration: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Name: "cortex_ruler_rule_group_sync_duration_seconds",
			Help: "The duration in seconds required to sync and load rule groups from storage.",
		}),
	}

	if len(cfg.EnabledTenants) > 0 {
		level.Info(ruler.logger).Log("msg", "ruler using enabled users", "enabled", strings.Join(cfg.EnabledTenants, ", "))
	}
	if len(cfg.DisabledTenants) > 0 {
		level.Info(ruler.logger).Log("msg", "ruler using disabled users", "disabled", strings.Join(cfg.DisabledTenants, ", "))
	}

	if cfg.EnableSharding {
		ringStore, err := kv.NewClient(
			cfg.Ring.KVStore,
			ring.GetCodec(),
			kv.RegistererWithKVName(prometheus.WrapRegistererWithPrefix("cortex_", reg), "ruler"),
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create KV store client")
		}

		if err = enableSharding(ruler, ringStore); err != nil {
			return nil, errors.Wrap(err, "setup ruler sharding ring")
		}
	}

	ruler.Service = services.NewBasicService(ruler.starting, ruler.run, ruler.stopping)
	return ruler, nil
}

func enableSharding(r *Ruler, ringStore kv.Client) error {
	lifecyclerCfg, err := r.cfg.Ring.ToLifecyclerConfig(r.logger)
	if err != nil {
		return errors.Wrap(err, "failed to initialize ruler's lifecycler config")
	}

	// Define lifecycler delegates in reverse order (last to be called defined first because they're
	// chained via "next delegate").
	delegate := ring.BasicLifecyclerDelegate(r)
	delegate = ring.NewLeaveOnStoppingDelegate(delegate, r.logger)
	delegate = ring.NewAutoForgetDelegate(r.cfg.Ring.HeartbeatTimeout*ringAutoForgetUnhealthyPeriods, delegate, r.logger)

	rulerRingName := "ruler"
	r.lifecycler, err = ring.NewBasicLifecycler(lifecyclerCfg, rulerRingName, ringKey, ringStore, delegate, r.logger, prometheus.WrapRegistererWithPrefix("cortex_", r.registry))
	if err != nil {
		return errors.Wrap(err, "failed to initialize ruler's lifecycler")
	}

	r.ring, err = ring.NewWithStoreClientAndStrategy(r.cfg.Ring.ToRingConfig(), rulerRingName, ringKey, ringStore, ring.NewIgnoreUnhealthyInstancesReplicationStrategy(), prometheus.WrapRegistererWithPrefix("cortex_", r.registry), r.logger)
	if err != nil {
		return errors.Wrap(err, "failed to initialize ruler's ring")
	}

	return nil
}

func (r *Ruler) starting(ctx context.Context) error {
	// If sharding is enabled, start the used subservices.
	if r.cfg.EnableSharding {
		var err error

		if r.subservices, err = services.NewManager(r.lifecycler, r.ring, r.clientsPool); err != nil {
			return errors.Wrap(err, "unable to start ruler subservices")
		}

		r.subservicesWatcher = services.NewFailureWatcher()
		r.subservicesWatcher.WatchManager(r.subservices)

		if err = services.StartManagerAndAwaitHealthy(ctx, r.subservices); err != nil {
			return errors.Wrap(err, "unable to start ruler subservices")
		}
	}

	// TODO: ideally, ruler would wait until its queryable is finished starting.
	return nil
}

// Stop stops the Ruler.
// Each function of the ruler is terminated before leaving the ring
func (r *Ruler) stopping(_ error) error {
	r.manager.Stop()

	if r.subservices != nil {
		_ = services.StopManagerAndAwaitStopped(context.Background(), r.subservices)
	}
	return nil
}

type sender interface {
	Send(alerts ...*notifier.Alert)
}

// SendAlerts implements a rules.NotifyFunc for a Notifier.
// It filters any non-firing alerts from the input.
//
// Copied from Prometheus's main.go.
func SendAlerts(n sender, externalURL string) promRules.NotifyFunc {
	return func(ctx context.Context, expr string, alerts ...*promRules.Alert) {
		var res []*notifier.Alert

		for _, alert := range alerts {
			a := &notifier.Alert{
				StartsAt:     alert.FiredAt,
				Labels:       alert.Labels,
				Annotations:  alert.Annotations,
				GeneratorURL: externalURL + strutil.TableLinkForExpression(expr),
			}
			if !alert.ResolvedAt.IsZero() {
				a.EndsAt = alert.ResolvedAt
			} else {
				a.EndsAt = alert.ValidUntil
			}
			res = append(res, a)
		}

		if len(alerts) > 0 {
			n.Send(res...)
		}
	}
}

func ruleGroupDisabled(ruleGroup *rulespb.RuleGroupDesc, disabledRuleGroupsForUser validation.DisabledRuleGroups) bool {
	for _, disabledRuleGroupForUser := range disabledRuleGroupsForUser {
		if ruleGroup.Namespace == disabledRuleGroupForUser.Namespace &&
			ruleGroup.Name == disabledRuleGroupForUser.Name &&
			ruleGroup.User == disabledRuleGroupForUser.User {
			return true
		}
	}
	return false
}

var sep = []byte("/")

func tokenForGroup(g *rulespb.RuleGroupDesc) uint32 {
	ringHasher := fnv.New32a()

	// Hasher never returns err.
	_, _ = ringHasher.Write([]byte(g.User))
	_, _ = ringHasher.Write(sep)
	_, _ = ringHasher.Write([]byte(g.Namespace))
	_, _ = ringHasher.Write(sep)
	_, _ = ringHasher.Write([]byte(g.Name))

	return ringHasher.Sum32()
}

func instanceOwnsRuleGroup(r ring.ReadRing, g *rulespb.RuleGroupDesc, disabledRuleGroups validation.DisabledRuleGroups, instanceAddr string) (bool, error) {

	hash := tokenForGroup(g)

	rlrs, err := r.Get(hash, RingOp, nil, nil, nil)
	if err != nil {
		return false, errors.Wrap(err, "error reading ring to verify rule group ownership")
	}

	ownsRuleGroup := false
	for _, instance := range rlrs.Instances {
		if instance.Addr == instanceAddr {
			ownsRuleGroup = true
			break
		}
	}

	if ownsRuleGroup && ruleGroupDisabled(g, disabledRuleGroups) {
		return false, &DisabledRuleGroupErr{Message: fmt.Sprintf("rule group %s, namespace %s, user %s is disabled", g.Name, g.Namespace, g.User)}
	}

	return ownsRuleGroup, nil
}

func (r *Ruler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.cfg.EnableSharding {
		r.ring.ServeHTTP(w, req)
	} else {
		var unshardedPage = `
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset="UTF-8">
					<title>Cortex Ruler Status</title>
				</head>
				<body>
					<h1>Cortex Ruler Status</h1>
					<p>Ruler running with shards disabled</p>
				</body>
			</html>`
		util.WriteHTMLResponse(w, unshardedPage)
	}
}

func (r *Ruler) run(ctx context.Context) error {
	level.Info(r.logger).Log("msg", "ruler up and running")

	tick := time.NewTicker(r.cfg.PollInterval)
	defer tick.Stop()

	var ringTickerChan <-chan time.Time
	var ringLastState ring.ReplicationSet

	if r.cfg.EnableSharding {
		ringLastState, _ = r.ring.GetAllHealthy(RingOp)
		ringTicker := time.NewTicker(util.DurationWithJitter(r.cfg.RingCheckPeriod, 0.2))
		defer ringTicker.Stop()
		ringTickerChan = ringTicker.C
	}

	r.syncRules(ctx, rulerSyncReasonInitial)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			r.syncRules(ctx, rulerSyncReasonPeriodic)
		case <-ringTickerChan:
			// We ignore the error because in case of error it will return an empty
			// replication set which we use to compare with the previous state.
			currRingState, _ := r.ring.GetAllHealthy(RingOp)

			if ring.HasReplicationSetChanged(ringLastState, currRingState) {
				ringLastState = currRingState
				r.syncRules(ctx, rulerSyncReasonRingChange)
			}
		case err := <-r.subservicesWatcher.Chan():
			return errors.Wrap(err, "ruler subservice failed")
		}
	}
}

func (r *Ruler) syncRules(ctx context.Context, reason string) {
	level.Debug(r.logger).Log("msg", "syncing rules", "reason", reason)
	r.rulerSync.WithLabelValues(reason).Inc()
	timer := prometheus.NewTimer(nil)

	defer func() {
		ruleGroupSyncDuration := timer.ObserveDuration().Seconds()
		r.ruleGroupSyncDuration.Set(ruleGroupSyncDuration)
	}()

	loadedConfigs, err := r.loadRuleGroups(ctx)
	if err != nil {
		return
	}

	// This will also delete local group files for users that are no longer in 'configs' map.
	r.manager.SyncRuleGroups(ctx, loadedConfigs)
}

func (r *Ruler) loadRuleGroups(ctx context.Context) (map[string]rulespb.RuleGroupList, error) {
	timer := prometheus.NewTimer(nil)

	defer func() {
		storeLoadSeconds := timer.ObserveDuration().Seconds()
		r.ruleGroupStoreLoadDuration.Set(storeLoadSeconds)
	}()

	configs, err := r.listRules(ctx)
	if err != nil {
		level.Error(r.logger).Log("msg", "unable to list rules", "err", err)
		return nil, err
	}

	loadedConfigs, err := r.store.LoadRuleGroups(ctx, configs)
	if err != nil {
		level.Warn(r.logger).Log("msg", "failed to load some rules owned by this ruler", "count", len(configs)-len(loadedConfigs), "err", err)
	}
	return loadedConfigs, nil
}

func (r *Ruler) listRules(ctx context.Context) (result map[string]rulespb.RuleGroupList, err error) {
	switch {
	case !r.cfg.EnableSharding:
		result, err = r.listRulesNoSharding(ctx)

	case r.cfg.ShardingStrategy == util.ShardingStrategyDefault:
		result, err = r.listRulesShardingDefault(ctx)

	case r.cfg.ShardingStrategy == util.ShardingStrategyShuffle:
		result, err = r.listRulesShuffleSharding(ctx)

	default:
		return nil, errors.New("invalid sharding configuration")
	}

	if err != nil {
		return
	}

	for userID := range result {
		if !r.allowedTenants.IsAllowed(userID) {
			level.Debug(r.logger).Log("msg", "ignoring rule groups for user, not allowed", "user", userID)
			delete(result, userID)
		}
	}
	return
}

func (r *Ruler) listRulesNoSharding(ctx context.Context) (map[string]rulespb.RuleGroupList, error) {
	allRuleGroups, err := r.store.ListAllRuleGroups(ctx)
	if err != nil {
		return nil, err
	}
	for userID, groups := range allRuleGroups {
		disabledRuleGroupsForUser := r.limits.DisabledRuleGroups(userID)
		if len(disabledRuleGroupsForUser) == 0 {
			continue
		}
		filteredGroupsForUser := rulespb.RuleGroupList{}
		for _, group := range groups {
			if !ruleGroupDisabled(group, disabledRuleGroupsForUser) {
				filteredGroupsForUser = append(filteredGroupsForUser, group)
			} else {
				level.Info(r.logger).Log("msg", "rule group disabled", "name", group.Name, "namespace", group.Namespace, "user", group.User)
			}
		}
		allRuleGroups[userID] = filteredGroupsForUser
	}
	return allRuleGroups, nil
}

func (r *Ruler) listRulesShardingDefault(ctx context.Context) (map[string]rulespb.RuleGroupList, error) {
	configs, err := r.store.ListAllRuleGroups(ctx)
	if err != nil {
		return nil, err
	}

	filteredConfigs := make(map[string]rulespb.RuleGroupList)
	for userID, groups := range configs {
		filtered := filterRuleGroups(userID, groups, r.limits.DisabledRuleGroups(userID), r.ring, r.lifecycler.GetInstanceAddr(), r.logger, r.ringCheckErrors)
		if len(filtered) > 0 {
			filteredConfigs[userID] = filtered
		}
	}
	return filteredConfigs, nil
}

func (r *Ruler) listRulesShuffleSharding(ctx context.Context) (map[string]rulespb.RuleGroupList, error) {
	users, err := r.store.ListAllUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list users of ruler")
	}

	// Only users in userRings will be used in the to load the rules.
	userRings := map[string]ring.ReadRing{}
	for _, u := range users {
		if shardSize := r.limits.RulerTenantShardSize(u); shardSize > 0 {
			subRing := r.ring.ShuffleShard(u, shardSize)

			// Include the user only if it belongs to this ruler shard.
			if subRing.HasInstance(r.lifecycler.GetInstanceID()) {
				userRings[u] = subRing
			}
		} else {
			// A shard size of 0 means shuffle sharding is disabled for this specific user.
			// In that case we use the full ring so that rule groups will be sharded across all rulers.
			userRings[u] = r.ring
		}
	}

	if len(userRings) == 0 {
		return nil, nil
	}

	userCh := make(chan string, len(userRings))
	for u := range userRings {
		userCh <- u
	}
	close(userCh)

	mu := sync.Mutex{}
	result := map[string]rulespb.RuleGroupList{}

	concurrency := loadRulesConcurrency
	if len(userRings) < concurrency {
		concurrency = len(userRings)
	}

	g, gctx := errgroup.WithContext(ctx)
	for i := 0; i < concurrency; i++ {
		g.Go(func() error {
			for userID := range userCh {
				groups, err := r.store.ListRuleGroupsForUserAndNamespace(gctx, userID, "")
				if err != nil {
					return errors.Wrapf(err, "failed to fetch rule groups for user %s", userID)
				}

				filtered := filterRuleGroups(userID, groups, r.limits.DisabledRuleGroups(userID), userRings[userID], r.lifecycler.GetInstanceAddr(), r.logger, r.ringCheckErrors)
				if len(filtered) == 0 {
					continue
				}

				mu.Lock()
				result[userID] = filtered
				mu.Unlock()
			}
			return nil
		})
	}

	err = g.Wait()
	return result, err
}

// filterRuleGroups returns map of rule groups that given instance "owns" based on supplied ring.
// This function only uses User, Namespace, and Name fields of individual RuleGroups.
//
// Reason why this function is not a method on Ruler is to make sure we don't accidentally use r.ring,
// but only ring passed as parameter.
func filterRuleGroups(userID string, ruleGroups []*rulespb.RuleGroupDesc, disabledRuleGroups validation.DisabledRuleGroups, ring ring.ReadRing, instanceAddr string, log log.Logger, ringCheckErrors prometheus.Counter) []*rulespb.RuleGroupDesc {
	// Prune the rule group to only contain rules that this ruler is responsible for, based on ring.
	var result []*rulespb.RuleGroupDesc
	for _, g := range ruleGroups {
		owned, err := instanceOwnsRuleGroup(ring, g, disabledRuleGroups, instanceAddr)
		if err != nil {
			switch e := err.(type) {
			case *DisabledRuleGroupErr:
				level.Info(log).Log("msg", e.Message)
				continue
			default:
				ringCheckErrors.Inc()
				level.Error(log).Log("msg", "failed to check if the ruler replica owns the rule group", "user", userID, "namespace", g.Namespace, "group", g.Name, "err", err)
				continue
			}
		}

		if owned {
			level.Debug(log).Log("msg", "rule group owned", "user", g.User, "namespace", g.Namespace, "name", g.Name)
			result = append(result, g)
		} else {
			level.Debug(log).Log("msg", "rule group not owned, ignoring", "user", g.User, "namespace", g.Namespace, "name", g.Name)
		}
	}

	return result
}

// GetRules retrieves the running rules from this ruler and all running rulers in the ring if
// sharding is enabled
func (r *Ruler) GetRules(ctx context.Context, rulesRequest RulesRequest) ([]*GroupStateDesc, error) {
	userID, err := tenant.TenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("no user id found in context")
	}

	if r.cfg.EnableSharding {
		return r.getShardedRules(ctx, userID, rulesRequest)
	}

	return r.getLocalRules(userID, rulesRequest)
}

func (r *Ruler) getLocalRules(userID string, rulesRequest RulesRequest) ([]*GroupStateDesc, error) {
	groups := r.manager.GetRules(userID)

	groupDescs := make([]*GroupStateDesc, 0, len(groups))
	prefix := filepath.Join(r.cfg.RulePath, userID) + "/"

	sliceToSet := func(values []string) map[string]struct{} {
		set := make(map[string]struct{}, len(values))
		for _, v := range values {
			set[v] = struct{}{}
		}
		return set
	}

	ruleNameSet := sliceToSet(rulesRequest.RuleNames)
	ruleGroupNameSet := sliceToSet(rulesRequest.RuleGroupNames)
	fileSet := sliceToSet(rulesRequest.Files)
	ruleType := rulesRequest.Type

	returnAlerts := ruleType == "" || ruleType == alertingRuleFilter
	returnRecording := ruleType == "" || ruleType == recordingRuleFilter

	for _, group := range groups {
		// The mapped filename is url path escaped encoded to make handling `/` characters easier
		decodedNamespace, err := url.PathUnescape(strings.TrimPrefix(group.File(), prefix))
		r.logger.Log("msg", "getLocalRules", "userID", userID, "decodedNamespace", decodedNamespace, "groupName", group.Name(), "ruleNameSet", strings.Join(rulesRequest.RuleNames, ";"), "ruleGroupNameSet", strings.Join(rulesRequest.RuleGroupNames, ";"), "fileSet", strings.Join(rulesRequest.Files, ";"))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode rule filename")
		}
		if len(fileSet) > 0 {
			if _, OK := fileSet[decodedNamespace]; !OK {
				continue
			}
		}

		if len(ruleGroupNameSet) > 0 {
			if _, OK := ruleGroupNameSet[group.Name()]; !OK {
				continue
			}
		}
		interval := group.Interval()

		groupDesc := &GroupStateDesc{
			Group: &rulespb.RuleGroupDesc{
				Name:      group.Name(),
				Namespace: string(decodedNamespace),
				Interval:  interval,
				User:      userID,
				Limit:     int64(group.Limit()),
			},

			EvaluationTimestamp: group.GetLastEvaluation(),
			EvaluationDuration:  group.GetEvaluationTime(),
		}
		for _, r := range group.Rules() {
			if len(ruleNameSet) > 0 {
				if _, OK := ruleNameSet[r.Name()]; !OK {
					continue
				}
			}
			lastError := ""
			if r.LastError() != nil {
				lastError = r.LastError().Error()
			}

			var ruleDesc *RuleStateDesc
			switch rule := r.(type) {
			case *promRules.AlertingRule:
				if !returnAlerts {
					continue
				}
				alerts := []*AlertStateDesc{}
				for _, a := range rule.ActiveAlerts() {
					alerts = append(alerts, &AlertStateDesc{
						State:       a.State.String(),
						Labels:      cortexpb.FromLabelsToLabelAdapters(a.Labels),
						Annotations: cortexpb.FromLabelsToLabelAdapters(a.Annotations),
						Value:       a.Value,
						ActiveAt:    a.ActiveAt,
						FiredAt:     a.FiredAt,
						ResolvedAt:  a.ResolvedAt,
						LastSentAt:  a.LastSentAt,
						ValidUntil:  a.ValidUntil,
					})
				}
				ruleDesc = &RuleStateDesc{
					Rule: &rulespb.RuleDesc{
						Expr:          rule.Query().String(),
						Alert:         rule.Name(),
						For:           rule.HoldDuration(),
						KeepFiringFor: rule.KeepFiringFor(),
						Labels:        cortexpb.FromLabelsToLabelAdapters(rule.Labels()),
						Annotations:   cortexpb.FromLabelsToLabelAdapters(rule.Annotations()),
					},
					State:               rule.State().String(),
					Health:              string(rule.Health()),
					LastError:           lastError,
					Alerts:              alerts,
					EvaluationTimestamp: rule.GetEvaluationTimestamp(),
					EvaluationDuration:  rule.GetEvaluationDuration(),
				}
			case *promRules.RecordingRule:
				if !returnRecording {
					continue
				}
				ruleDesc = &RuleStateDesc{
					Rule: &rulespb.RuleDesc{
						Record: rule.Name(),
						Expr:   rule.Query().String(),
						Labels: cortexpb.FromLabelsToLabelAdapters(rule.Labels()),
					},
					Health:              string(rule.Health()),
					LastError:           lastError,
					EvaluationTimestamp: rule.GetEvaluationTimestamp(),
					EvaluationDuration:  rule.GetEvaluationDuration(),
				}
			default:
				return nil, errors.Errorf("failed to assert type of rule '%v'", rule.Name())
			}
			groupDesc.ActiveRules = append(groupDesc.ActiveRules, ruleDesc)
		}
		if len(groupDesc.ActiveRules) > 0 {
			groupDescs = append(groupDescs, groupDesc)
		}
	}

	return groupDescs, nil
}

func (r *Ruler) getShardedRules(ctx context.Context, userID string, rulesRequest RulesRequest) ([]*GroupStateDesc, error) {
	ring := ring.ReadRing(r.ring)

	if shardSize := r.limits.RulerTenantShardSize(userID); shardSize > 0 && r.cfg.ShardingStrategy == util.ShardingStrategyShuffle {
		ring = r.ring.ShuffleShard(userID, shardSize)
	}

	rulers, err := ring.GetReplicationSetForOperation(RingOp)
	if err != nil {
		return nil, err
	}

	ctx, err = user.InjectIntoGRPCRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to inject user ID into grpc request, %v", err)
	}

	var (
		mergedMx sync.Mutex
		merged   []*GroupStateDesc
	)

	// Concurrently fetch rules from all rulers. Since rules are not replicated,
	// we need all requests to succeed.
	jobs := concurrency.CreateJobsFromStrings(rulers.GetAddresses())
	err = concurrency.ForEach(ctx, jobs, len(jobs), func(ctx context.Context, job interface{}) error {
		addr := job.(string)

		rulerClient, err := r.clientsPool.GetClientFor(addr)
		if err != nil {
			return errors.Wrapf(err, "unable to get client for ruler %s", addr)
		}

		newGrps, err := rulerClient.Rules(ctx, &RulesRequest{
			RuleNames:      rulesRequest.GetRuleNames(),
			RuleGroupNames: rulesRequest.GetRuleGroupNames(),
			Files:          rulesRequest.GetFiles(),
			Type:           rulesRequest.GetType(),
		})
		if err != nil {
			return errors.Wrapf(err, "unable to retrieve rules from ruler %s", addr)
		}

		mergedMx.Lock()
		merged = append(merged, newGrps.Groups...)
		mergedMx.Unlock()

		return nil
	})

	return merged, err
}

// Rules implements the rules service
func (r *Ruler) Rules(ctx context.Context, in *RulesRequest) (*RulesResponse, error) {
	userID, err := tenant.TenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("no user id found in context")
	}

	groupDescs, err := r.getLocalRules(userID, *in)
	if err != nil {
		return nil, err
	}

	return &RulesResponse{Groups: groupDescs}, nil
}

// AssertMaxRuleGroups limit has not been reached compared to the current
// number of total rule groups in input and returns an error if so.
func (r *Ruler) AssertMaxRuleGroups(userID string, rg int) error {
	limit := r.limits.RulerMaxRuleGroupsPerTenant(userID)

	if limit <= 0 {
		return nil
	}

	if rg <= limit {
		return nil
	}

	return fmt.Errorf(errMaxRuleGroupsPerUserLimitExceeded, limit, rg)
}

// AssertMaxRulesPerRuleGroup limit has not been reached compared to the current
// number of rules in a rule group in input and returns an error if so.
func (r *Ruler) AssertMaxRulesPerRuleGroup(userID string, rules int) error {
	limit := r.limits.RulerMaxRulesPerRuleGroup(userID)

	if limit <= 0 {
		return nil
	}

	if rules <= limit {
		return nil
	}
	return fmt.Errorf(errMaxRulesPerRuleGroupPerUserLimitExceeded, limit, rules)
}

func (r *Ruler) DeleteTenantConfiguration(w http.ResponseWriter, req *http.Request) {
	logger := util_log.WithContext(req.Context(), r.logger)

	userID, err := tenant.TenantID(req.Context())
	if err != nil {
		// When Cortex is running, it uses Auth Middleware for checking X-Scope-OrgID and injecting tenant into context.
		// Auth Middleware sends http.StatusUnauthorized if X-Scope-OrgID is missing, so we do too here, for consistency.
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = r.store.DeleteNamespace(req.Context(), userID, "") // Empty namespace = delete all rule groups.
	if err != nil && !errors.Is(err, rulestore.ErrGroupNamespaceNotFound) {
		respondError(logger, w, err.Error())
		return
	}

	level.Info(logger).Log("msg", "deleted all tenant rule groups", "user", userID)
	w.WriteHeader(http.StatusOK)
}

func (r *Ruler) ListAllRules(w http.ResponseWriter, req *http.Request) {
	logger := util_log.WithContext(req.Context(), r.logger)

	userIDs, err := r.store.ListAllUsers(req.Context())
	if err != nil {
		level.Error(logger).Log("msg", errListAllUser, "err", err)
		http.Error(w, fmt.Sprintf("%s: %s", errListAllUser, err.Error()), http.StatusInternalServerError)
		return
	}

	done := make(chan struct{})
	iter := make(chan interface{})

	go func() {
		util.StreamWriteYAMLV3Response(w, iter, logger)
		close(done)
	}()

	err = concurrency.ForEachUser(req.Context(), userIDs, fetchRulesConcurrency, func(ctx context.Context, userID string) error {
		rg, err := r.store.ListRuleGroupsForUserAndNamespace(ctx, userID, "")
		if err != nil {
			return errors.Wrapf(err, "failed to fetch ruler config for user %s", userID)
		}
		userRules := map[string]rulespb.RuleGroupList{userID: rg}
		if userRules, err = r.store.LoadRuleGroups(ctx, userRules); err != nil {
			return errors.Wrapf(err, "failed to load ruler config for user %s", userID)
		}
		data := map[string]map[string][]rulefmt.RuleGroup{userID: userRules[userID].Formatted()}

		select {
		case iter <- data:
		case <-done: // stop early, if sending response has already finished
		}

		return nil
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to list all ruler configs", "err", err)
	}
	close(iter)
	<-done
}
