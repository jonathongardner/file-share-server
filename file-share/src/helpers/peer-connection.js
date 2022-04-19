export default class PeerConnection {
  constructor(identifier) {
    const pc = new RTCPeerConnection({ sdpSemantics: 'unified-plan' })

    // connect audio / video
    pc.addEventListener('icegatheringstatechange', this.iceStateChange)
    pc.addEventListener('connectionstatechange', () => {
      this.stateCallback(this.identifier, this.pc.connectionState)
    });
    pc.addEventListener('datachannel', ({ channel }) => {
      if (channel.label === 'file' && !this.fileC) {
        this.fileC = channel
        this.fileC.addEventListener('message', (event) => {
          this.fileCallback(this.identifier, JSON.parse(event.data))
        })
      }
    });

    this.identifier = identifier
    this.pc = pc
  }

  iceStateChange = () => {
    const { iceGatheringState } = this.pc
    console.log(iceGatheringState)

    if (iceGatheringState === 'complete') {
      this.descriptionCallback(this.identifier, this.pc.localDescription)
    }
  }

  createLocalDescription = async () => {
    console.log('Creating local description')
    // if we are sending request than we should have data channel
    this.fileC = this.pc.createDataChannel('file')
    this.fileC.addEventListener('message', (event) => {
      this.fileCallback(this.identifier, JSON.parse(event.data))
    })

    const offer = await this.pc.createOffer()
    // after offer is set it will call iceStateChange
    await this.pc.setLocalDescription(offer)
  }
  setRemoteDescription = async (description) => {
    console.log('Setting remote description')

    await this.pc.setRemoteDescription(description)
    if (!this.pc.localDescription) { // || description.type === 'offer'
      console.log('Creating local description')
      const des = await this.pc.createAnswer()
      // after offer is set it will call iceStateChange
      await this.pc.setLocalDescription(des)
    }
  }
  stop = async () => {
    setTimeout(() => {
      this.pc.close()
    }, 500)
  }
  // override
  stateCallback = () => null
  descriptionCallback = () => null
  fileCallback = () => null
  // closeCallback = () => { console.log('closed') }
}
