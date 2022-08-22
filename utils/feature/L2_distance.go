package feature

// FeatureVector 特征向量
type FeatureVector []float32

// Sub 计算2个特征向量的差向量
func (fv FeatureVector) sub(targetFv FeatureVector) FeatureVector {
	resFv := make(FeatureVector, len(fv))
	for k := range fv {
		resFv[k] = fv[k] - targetFv[k]
	}
	return resFv
}

// 计算特征向量的平方和
func (fv FeatureVector) squareSum() float64 {
	var squareSum float64
	for k := range fv {
		squareSum += float64(fv[k] * fv[k])
	}
	return squareSum
}

// 计算2个特征向量的L2距离
func (fv FeatureVector) getL2Distance(targetFv FeatureVector) float64 {
	return fv.sub(targetFv).squareSum()
}

// FeatureData 特征base64字符串
type FeatureData string
