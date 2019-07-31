package forwarder

import (
	"autoscaler/metricsforwarder/config"
	"autoscaler/models"

	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/lager"
)

type EnvelopeStreamerLogger struct {
	logger lager.Logger
}

func (l *EnvelopeStreamerLogger) Printf(message string, data ...interface{}) {
	l.logger.Debug(message, lager.Data{"data": data})
}
func (l *EnvelopeStreamerLogger) Panicf(message string, data ...interface{}) {
	l.logger.Fatal(message, nil, lager.Data{"data": data})
}

type MetricForwarder interface {
	EmitMetric(*models.CustomMetric)
}

type metricForwarder struct {
	client *loggregator.IngressClient
	logger lager.Logger
}

const METRICS_FORWARDER_ORIGIN = "autoscaler_metrics_forwarder"

func NewMetricForwarder(logger lager.Logger, conf *config.Config) (MetricForwarder, error) {
	tlsConfig, err := loggregator.NewIngressTLSConfig(
		conf.LoggregatorConfig.TLS.CACertFile,
		conf.LoggregatorConfig.TLS.CertFile,
		conf.LoggregatorConfig.TLS.KeyFile,
	)
	if err != nil {
		logger.Error("could-not-create-TLS-config", err, lager.Data{"config": conf})
		return &metricForwarder{}, err
	}

	client, err := loggregator.NewIngressClient(
		tlsConfig,
		loggregator.WithAddr(conf.LoggregatorConfig.MetronAddress),
		loggregator.WithTag("origin", METRICS_FORWARDER_ORIGIN),
		loggregator.WithLogger(&EnvelopeStreamerLogger{
			logger: logger.Session("emitter"),
		}),
	)

	if err != nil {
		logger.Error("could-not-create-loggregator-client", err, lager.Data{"config": conf})
		return &metricForwarder{}, err
	}

	return &metricForwarder{
		client: client,
		logger: logger,
	}, nil
}

func (mf *metricForwarder) EmitMetric(metric *models.CustomMetric) {
	mf.logger.Debug("custom-metric-emit-request-received:", lager.Data{"metric": metric})

	options := []loggregator.EmitGaugeOption{
		loggregator.WithGaugeAppInfo(metric.AppGUID, int(metric.InstanceIndex)),
		loggregator.WithGaugeValue(metric.Name, metric.Value, metric.Unit),
	}
	mf.client.EmitGauge(options...)
}
