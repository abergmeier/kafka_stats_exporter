package types

type LabelNameTransformer func(value string) (labelName string)
type MetricNameTransformer func(value string) (labelName string)
