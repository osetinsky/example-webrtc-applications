package jack

import (
  "flag"
  "fmt"
  "math/rand"

  "github.com/pion/webrtc/v2"

  gst "github.com/osetinsky/example-webrtc-applications/internal/gstreamer-src"
  "github.com/osetinsky/example-webrtc-applications/internal/signal"
)

func StartGstreamer(flagName, browserToken string, ch chan string) {
  go func() {

    // gst := fmt.Sprintf(`echo %s | gstreamer-send -audio-src "jackaudiosrc ! audioconvert ! audioresample"`, t.Test)
    // audioSrc := flag.String("audio-src", "jackaudiosrc ! audioconvert ! audioresample", "GStreamer audio src")
    audioSrc := flag.String(flagName, "jackaudiosrc ! audioconvert ! audioresample", "GStreamer audio src")

    flag.Parse()

    // Everything below is the pion-WebRTC API! Thanks for using it ❤️.

    // Prepare the configuration
    config := webrtc.Configuration{
      ICEServers: []webrtc.ICEServer{
        {
          URLs: []string{"stun:stun.l.google.com:19302"},
        },
      },
    }

    // Create a new RTCPeerConnection
    peerConnection, err := webrtc.NewPeerConnection(config)
    if err != nil {
      panic(err)
    }

    // Set the handler for ICE connection state
    // This will notify you when the peer has connected/disconnected
    peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
      // fmt.Printf("Connection State has changed %s \n", connectionState.String())
    })

    // Create a audio track
    audioTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeOpus, rand.Uint32(), "audio", "pion1")
    if err != nil {
      panic(err)
    }
    _, err = peerConnection.AddTrack(audioTrack)
    if err != nil {
      panic(err)
    }

    // Wait for the offer to be pasted
    offer := webrtc.SessionDescription{}
    // signal.Decode(signal.MustReadStdin(), &offer)
    signal.Decode(browserToken, &offer)

    // Set the remote SessionDescription
    err = peerConnection.SetRemoteDescription(offer)
    if err != nil {
      panic(err)
    }

    // Create an answer
    answer, err := peerConnection.CreateAnswer(nil)
    if err != nil {
      panic(err)
    }

    // Sets the LocalDescription, and starts our UDP listeners
    err = peerConnection.SetLocalDescription(answer)
    if err != nil {
      panic(err)
    }

    ch <-signal.Encode(answer)
    // close(ch)

    gst.CreatePipeline(webrtc.Opus, []*webrtc.Track{audioTrack}, *audioSrc).Start()

    // Block forever
    select {}
  }()
}
