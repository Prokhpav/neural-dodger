package run

import (
	"math"
	"math/rand"
)

const (
	InpCount = 20
	MidCount = 5
	OutCount = 1
)

var (
	MuteAmplitude = 2.
	MuteChance    = 15.
	AllNetworks   []Network
)

func NeuronFunc(x float64) float64 {
	return x/(1+math.Abs(x))/2 + 0.5
}

type Network struct {
	SwInpMid [InpCount + 1][MidCount]float64
	SwMidOut [MidCount + 1][OutCount]float64
	Fitness  float64
}

func (n Network) Get(inpInp [InpCount]float64) [OutCount]float64 {
	var inpOut [InpCount + 1]float64
	var mid [MidCount + 1]float64
	var out [OutCount]float64

	for i := 0; i < len(inpInp); i++ {
		inpOut[i] = NeuronFunc(inpInp[i])
	}
	inpOut[InpCount] = 1
	for midI := 0; midI < MidCount; midI++ {
		for inpI := 0; inpI <= InpCount; inpI++ {
			mid[midI] += inpOut[inpI] * n.SwInpMid[inpI][midI]
		}
		mid[midI] = NeuronFunc(mid[midI])
	}
	mid[MidCount] = 1
	for outI := 0; outI < OutCount; outI++ {
		for midI := 0; midI <= MidCount; midI++ {
			out[outI] += mid[midI] * n.SwMidOut[midI][outI]
		}
		out[outI] = NeuronFunc(out[outI])
	}
	return out
}

func (n *Network) Mutate() { // Mutate one of the connections
	ind := rand.Intn((InpCount+1)*MidCount + (MidCount+1)*OutCount)
	if ind < (InpCount+1)*MidCount {
		n.SwInpMid[ind%(InpCount+1)][ind/(InpCount+1)] += (rand.Float64()*2 - 1) * MuteAmplitude
	} else {
		ind -= (InpCount + 1) * MidCount
		n.SwMidOut[ind%(MidCount+1)][ind/(MidCount+1)] += (rand.Float64()*2 - 1) * MuteAmplitude
	}
}

func GetRandomNetwork() (network Network) {
	for inpI := 0; inpI <= InpCount; inpI++ {
		for midI := 0; midI < MidCount; midI++ {
			network.SwInpMid[inpI][midI] = rand.Float64()*2 - 1
		}
	}
	for midI := 0; midI <= MidCount; midI++ {
		for outI := 0; outI < OutCount; outI++ {
			network.SwMidOut[midI][outI] = rand.Float64()*2 - 1
		}
	}
	network.Fitness = 0
	return
}

func GetOneOfNetworks() Network {
	kFunc := func(fitness float64) (chance float64) {
		if fitness > 0 {
			return math.Pow(fitness, 0.5) + 5
		}
		return 0
	}

	sumK := 0.
	for _, n := range AllNetworks {
		sumK += kFunc(n.Fitness)
	}
	r := rand.Float64() * sumK
	i := 0
	for r > 0 && i < len(AllNetworks) {
		r -= kFunc(AllNetworks[i].Fitness)
		i++
	}
	if i == 0 {
		return GetRandomNetwork()
	}
	return AllNetworks[i-1]
}

func NewNetwork() (netN Network) {
	net1 := GetOneOfNetworks()
	net2 := GetOneOfNetworks()
	for inpI := 0; inpI <= InpCount; inpI++ {
		for midI := 0; midI < MidCount; midI++ {
			if rand.Float64() < 0.5 {
				netN.SwInpMid[inpI][midI] = net1.SwInpMid[inpI][midI]
			} else { //                        |
				netN.SwInpMid[inpI][midI] = net2.SwInpMid[inpI][midI]
			}
		}
	}
	for midI := 0; midI <= MidCount; midI++ {
		for outI := 0; outI < OutCount; outI++ {
			if rand.Float64() < 0.5 {
				netN.SwMidOut[midI][outI] = net1.SwMidOut[midI][outI]
			} else { //                        |
				netN.SwMidOut[midI][outI] = net2.SwMidOut[midI][outI]
			}
		}
	}
	for muteCount := rand.Float64() * MuteChance; rand.Float64() < muteCount; muteCount-- {
		netN.Mutate()
	}
	return
}
