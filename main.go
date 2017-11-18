package main

import (
	"github.com/timshannon/go-openal"
	"fmt"
	"io/ioutil"
	"time"
	"math/cmplx"
	"github.com/mjibson/go-dsp/fft"
    //"math/rand"
	//"math"
	//"math"
	//"github.com/go-gl/gl/v2.1/gl"
	//"math"
	//"bytes"
	//"math"
	//"math"
)

func iterate(c complex128, iter int) complex128 {
	z := complex(0, 0)

	for i := 1; i < iter; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2000 {
			return z
		}
	}

	return z
}

func partitionfft ( wavfloat []float64 ) []complex128 {
	print(len(wavfloat))

	k:=0
	var fftresult []complex128 = make([]complex128,len(wavfloat))
	var win []float64 = make([]float64, 8000)
	var winfft []complex128 = make([]complex128,8000)
	for k<len(wavfloat) {
		j:=0

		for j<len(win) {
			win[j]=0
			j+=1
		}
		c:=0
		for k<len(wavfloat)&&c<len(win){
			win[c] = wavfloat[k]
			k+=1
			c+=1
		}

		winfft = fft.FFTReal(win)
		i:=0
		for m:=k-len(win); m<k; m++ {
			if m>len(wavfloat) {
				return fftresult
			}
			fftresult[m] = winfft[i]

			i+=1

		}

	}

  return fftresult
}

func blocktrash (wavfloat []float64) []float64 {

	count := 3
	blocklen := len(wavfloat)/count


	 var trash [][] float64 = make ([][]float64,count)
	 i:=0
	 for i<count{
	 	trash[i] = make ([] float64,blocklen)
	 	i+=1
	 }

	 i=0

	 k:=0
	 for i<count{
	 	j:=0
	 	for j<blocklen  {
	 		trash[i][j]=wavfloat[k]
	 		j+=1
	 		k+=1
	 		if k>len(wavfloat){
	 			break;
			}
		}
		i+=1
	 }

	 var sample []float64 = make([]float64,blocklen)
	 sample = trash[0]
	 trash[0]=trash[1]
	 trash[1]=sample


	result := []float64{}
	i=0
	for i<count{
		j:=0
		for j<blocklen {
		 result=append(result, trash[i][j])
		 j+=1
		}
		i+=1
	}

	 return result


}


func partitionifft (wavcomplex []complex128 ) []complex128  {

	k:=0
	var ifftresult []complex128 = make([]complex128,len(wavcomplex))
	var win []complex128 = make([]complex128, 2000)
	var winifft []complex128 = make([]complex128,2000)
	for k<len(wavcomplex) {
		j:=0

		for j<len(win) {
			win[j]=0
			j+=1
		}
		c:=0
		for k<len(wavcomplex)&&c<len(win){
			win[c] = wavcomplex[k]
			k+=1
			c+=1
		}

		winifft = fft.IFFT(win)
		i:=0
		for m:=k-len(win); m<k; m++ {
			if m>len(wavcomplex) {
				return ifftresult
			}
			ifftresult[m] = winifft[i]

			i+=1

		}

	}

	return ifftresult

}


func wavebytestofloat(wavBytes1 []byte) []float64 {
	var wavfloat1 []float64 = make([]float64, len(wavBytes1)/2)
	i := 0
	for i < len(wavBytes1)-1 {
		wavfloat1[i/2] = float64(int16(wavBytes1[i]) + int16(wavBytes1[i+1])<<8)
		i += 2
	}

	return wavfloat1

}

func wavefloattobytes(realpart []float64) []byte {
	i := 0
	var newWavBytes []byte = make([]byte, len(realpart)*2)
	for i < len(newWavBytes)-1 {
		test1 := byte(0xff & uint16(realpart[i/2]))
		test2 := byte((0xff00 & uint16(realpart[i/2])) >> 8)
		newWavBytes[i] = test1
		newWavBytes[i+1] = test2
		i += 2
	}
	return newWavBytes
}

func main() {

	// sets up OpenAL with default options
	device := openal.OpenDevice("")
	defer device.CloseDevice()
	context := device.CreateContext()
	defer context.Destroy()
	context.Activate()
	vendor := openal.GetVendor()
	if err := openal.Err(); err != nil {
		fmt.Printf("Failed to setup OpenAL: %v\n", err)
		return
	}
	fmt.Printf("OpenAL vendor: %s\n", vendor)
	for true {

	// make sure things have gone well

	source1 := openal.NewSource()
	defer source1.Pause()
	source1.SetLooping(false)
	source1.SetPosition(&openal.Vector{0.0, 0.0, -5})

	soundBuffer1 := openal.NewBuffer()
	soundBuffer2 := openal.NewBuffer()
	source2 := openal.NewSource()
	defer source2.Pause()
	source2.SetLooping(false)
	source2.SetPosition(&openal.Vector{0.0, 0.0, -5})

	if err := openal.Err(); err != nil {
		fmt.Printf("OpenAL buffer creation failed: %v\n", err)
		return
	}

	// load a sound effect
	wavBytes1, err := ioutil.ReadFile("assets/04 - Easy.wav")
	//wavBytes2, err := ioutil.ReadFile("assets/cartoon001.wav")

	//fft1 := fft.FFTReal(wavebytestofloat(wavBytes1))
	fft1 :=partitionfft(blocktrash(wavebytestofloat(wavBytes1)))
	//print(fft1)
	i:=0

	//r:=0
	for i<len(fft1){
		//r = rand.Intn(100)
		fft1[i] =  fft1[i] *(iterate(fft1[i],i)/(10e8))
		i+=1
	}



	ifft1 := partitionifft(fft1)

	var realpart = make([]float64, len(ifft1))

	i= 0

	for i < len(ifft1) {
		realpart[i] = real(ifft1[i])
		i += 1
	}

	var newWavBytes = wavefloattobytes(realpart)

	if err != nil {
		fmt.Printf("Failed to load the sound effect: %v\n", err)
		return
	}

	soundBuffer1.SetData(openal.FormatMono16, newWavBytes, 44100)
	source1.SetBuffer(soundBuffer1)

	soundBuffer2.SetData(openal.FormatMono16, wavBytes1, 44100)
	source2.SetBuffer(soundBuffer2)

	// play the sound
	source1.SetGain(5000)
	source1.Play()

	//source2.Play()

	for source1.State() == openal.Playing {
		// loop long enough to let the wave file finish
		time.Sleep(time.Millisecond * 0)
	}

	for source2.State() == openal.Playing {
		// loop long enough to let the wave file finish
		time.Sleep(time.Millisecond * 0)
	}

	source1.Stop()
	source2.Stop()

	//source1.Delete()
	//source2.Delete()
	}

	fmt.Println("Sound played!")
}
