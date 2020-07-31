package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v2"

	"github.com/sikang99/pion-radio-example/internal/signal"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	sdpInChan, sdpOutChan, rport := signal.HTTPSDPServer()

	m := webrtc.MediaEngine{}
	m.RegisterDefaultCodecs()

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	rtcConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pubConn, err := api.NewPeerConnection(rtcConfig)
	if err != nil {
		log.Fatalln(err)
	}

	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", rport))
	if err != nil {
		log.Println(err)
		return
	}

	udp, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer udp.Close()

	inRTPPacket := make([]byte, 4096) // UDP MTU size = 4096
	n, _, err := udp.ReadFromUDP(inRTPPacket)
	if err != nil {
		log.Fatalln(err)
	}

	packet := &rtp.Packet{}
	err = packet.Unmarshal(inRTPPacket[:n])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("SSRC:", packet.SSRC)

	localTrack, err := pubConn.NewTrack(webrtc.DefaultPayloadTypePCMU, packet.SSRC, "audio", "spider")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = pubConn.AddTrack(localTrack)
	if err != nil {
		log.Fatalln(err)
	}

	go RecvDataToTrack(udp, localTrack)

	for {
		log.Println("Now waiting for audio player connections ...")

		recvOnlyOffer := webrtc.SessionDescription{}
		signal.Decode(<-sdpInChan, &recvOnlyOffer)

		subConn, err := api.NewPeerConnection(rtcConfig)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = subConn.AddTrack(localTrack)
		if err != nil {
			log.Println(err)
			continue
		}

		err = subConn.SetRemoteDescription(recvOnlyOffer)
		if err != nil {
			log.Println(err)
			continue
		}

		answer, err := subConn.CreateAnswer(nil)
		if err != nil {
			log.Println(err)
			continue
		}

		err = subConn.SetLocalDescription(answer)
		if err != nil {
			log.Println(err)
			continue
		}

		playAnswer := signal.Encode(answer)
		sdpOutChan <- playAnswer
		//log.Println(playAnswer)
	}
}

func RecvDataToTrack(udp *net.UDPConn, toTrack *webrtc.Track) {
	inRTPPacket := make([]byte, 4096)

	for {
		n, _, err := udp.ReadFrom(inRTPPacket)
		if err != nil {
			log.Println(err)
			break
		}

		packet := &rtp.Packet{}
		err = packet.Unmarshal(inRTPPacket[:n])
		if err != nil {
			log.Println(err)
			break
		}

		err = toTrack.WriteRTP(packet)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
