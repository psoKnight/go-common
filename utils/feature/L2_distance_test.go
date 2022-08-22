package feature

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestGetFeatureL2Distance(t *testing.T) {
	a := "HekBCRgS9bMF7usdGgBA9Erv+voWHt3qAQkf3t303izGE8nj+tju5fIL/goGDA4pASkeFhToEdhD3hkuEyT8EtLm2ATuFvoA+NwaNA8sv/bc1wUFFf/lLs/1Nvn1DfOm+wP2AREbCMsNGvEf+w79yBXk9/Mu6/8H5yDz+h0YMNK+EPDZ2/v7Fw/j+dHkJxzzGBbf3wv2Dg4CEuo4zPjp1fPYCyPgBeK6+vgF8w//H/rn7wgc/d0mIgcGFfTu2QD9+szv7QsUEDHI6Aft/vII7+zxMLMCwhjL8Q4Z9vj8uyTy/PX7NBr9Crr5/dQW3eDy88IM8hHsHu3zBQ4Kyg0X9Q=="
	b := "OxbuPuLy7yv16efnBxIeG/30By3xDSng5PoN3BAt4f7/I9MR8uMX0fH7FCLz7QAlAvcHDPsSNPHp+Qj37R81GCkX3OsKEgLt/BsOBfdCSQoY2TDoBQ4X7wTyEu7nAwzS6RUCBSQPFeso9ynh6sXZCBIr6Efo8wEbzQL6GvjOHgvyZgrb7A3VzOmh8vjt2fYbBgq7+9oN/yzx0OA/HQIFy/8q9AbuBQju/vQYDOkZBdsF8isMJeblFfzg6fQGDvny9tD50iQG5R/wJOAQ/qkKF9/06ur2+fjYBO7p7Rz3Eir0JAUIASK37QwCzwcVA/H17tYXxznBPer56v8BAuPa9A=="

	// 同组
	//a := "INUd5x0x9eD3ufP6GeJLBfwOEyvAG9wGABQIy/tR9T0B/8rrFOnwyOvaDwAcDBUBAPEIAtf1B/o5/g72DiX7FfcF69weGOq18r4HBws++/LoB+QgE/sMEe09F+YC2/Xb79zXGQ71HfMw8gDy3hYFzA4ANwHyGuThzwMYCvkP6u7qKwT8Ehj/QkLT+/Hu8xnbBwgEtCLdF/ga2ulDKvfxy/vz9v/gESD00iscCy30EwcIBz4XLeUeKxr57fLZI/wMAtTN1g0b1A/o6CcnDNoK/BoXFeL21fbMCugw3svk5/sNFN759gz/0///1gAS6OYGHfr8GhX73zEB7BER+Bn2EA=="
	//b := "Ed8gCRAu59j1tRkEBPU1ChkDDiLTEPfr9gUNzfRQ/xbrD7jjGvDexADkBwkdJxIF9fwCF9fv8PovAxoDHg388vPs5OMxIOum7dvx9Ag8BuzIFd8PJ/wPBfJEFufr3unc39bUGyPsEPEF6gnj4wH+uhL5Lw7rIObt7/3/8/AV9P3MBwUEDPj/QzLN6+z49hTgBxDpwifaGv8L4/hN/gLavQbhBzPaERnf1S8d/iP2KAACDiYLKfoxBw/73QbeBgj6/dXhtgEUtgng9gokHNsf/QsXDtv4wgbS9Poi/tvg6gEeC+jh/iv93QYH5xEF3vIKFQX8KRwM1yzt5B4fABj7Cg=="

	// 同组
	//a := "NAT4GggL5PgX5doA/9UEFff83gITODcKydr7tQDZDOAPEPXz4QIa4e7sEO8LAUfsCSwlAA37Gl8zEzT2NzkqFe7x6rsZAOojG+QkLQQg2gAMruXR9SD81u7b5foA8sIDCPT49AfZLPblMBoJ4vMSEjr2/x322Q3s9hALCPUk5e3dPxnoN+71Au/K6Oky+AwK+QrMDsryByTX+gAs7f785+EI+OQB5RAK7wMN3/kELAEP+eLnLc8F4zwPFSHoBA/89ObgHPYV6wsSLOfq4a4p//TP8dAa4Ajy980aBQnwAS7v7gP4/BjYAvMAHfAWKh8V1NQJGhvfCfDU8DPsLf/16Q=="
	//b := "JfjXBOQaxO48Bt7iJuMOFPoI8TYaQy8B8/sIxfLy9A8cDwbcBAXj2f8YBeo1FhLtMRYy/Ofv/i8KD1z3DQPsBvr24gEDCPL/D+b2EPIoIvMdwfj47hUA4/vk/OYG+/jUHQTt/Or1MPrjMur76/r06Tjn1Trr/yO82QIiCvQh98nLBf7eNyz7GuqpC9gzAP4L9x3N/NDg8zLl++sK/PWq+9T7HQvT3gAH/w0Q3doKOgww4hvxHPUf4ULwOALjBgUR/tDcFw/3Bg4nDB7U5s8VEA3j+OgCBeXx5LcjLi/5FTXv/hnnCwMB8vUyIw0UIgwW/+jzAtT/FfXx+BMIFeHeBQ=="
	featureA := FeatureData(a)
	featureB := FeatureData(b)

	decodeFeatureA, err := featureA.getDecodeFeature()
	fmt.Println("A", len(decodeFeatureA))
	if err != nil {
		t.Errorf("decode feature A err: %v", err)
	}
	decodeFeatureB, err := featureB.getDecodeFeature()
	fmt.Println("B", len(decodeFeatureA))

	if err != nil {
		t.Errorf("decode feature B err: %v", err)
	}

	distance := decodeFeatureA.getL2Distance(decodeFeatureB)

	t.Log(distance)
}

// GetDecodeFeature 特征解码
func (f FeatureData) getDecodeFeature() (FeatureVector, error) {

	sF := string(f)
	fByte, err := base64.StdEncoding.DecodeString(sF)
	if err != nil {
		return nil, err
	}
	res := make([]float32, 0)
	for i := range fByte {
		// int8(bytes[i])*vdiff/scale(默认：255)+vmin
		res = append(res, float32(int8(fByte[i]))*8.662513732910156/255+(-4.304102420806885))
	}
	return FeatureVector(res), nil
}
