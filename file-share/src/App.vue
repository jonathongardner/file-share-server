<template>
  <h1>Clients</h1>
  <div v-for="pClient in pClients" :key='pClient.identifier' class='card'>
    <card-header :pClient='pClient' />
    <card-footer :pClient='pClient' @accept='acceptConnection' @decline='declineConnection' @disconnect='disconnect' @connect='connect'/>
  </div>
</template>

<script>
import CardHeader from './components/card-header.vue'
import CardFooter from './components/card-footer.vue'
import PeerConnection from './helpers/peer-connection.js'

export default {
  name: 'App',
  components: {
    CardHeader, CardFooter
  },
  data () {
    const host = process.env.VUE_SOCKET || location.host
    const socket = new WebSocket(`ws://${host}/candidates`);
    // socket.onopen = this.socketOpen
    socket.onmessage = this.socketMessge

    return {
      socket,
      clients: [{ identifier: 'a', name: 'bable', ipAddress: '1.2.3.4' }],
      peers: { 'a': { 'connectionRequest': 'poop' }}
    }
  },
  computed: {
    pClients () {
      return this.clients.map(c => {
        const peer = this.peers[c.identifier] || {}
        return {
          ...c,
          connectionRequest: !!peer.connectionRequest,
          connecting: !peer.connected && !!peer.pc,
          connected: !!peer.connected
        }
      })
    }
  },
  methods: {
    //-----Socket-----
    // socketOpen (message) {
    //   console.log("Socket opened")
    // },
    socketMessge (message) {
      const { type, data } = JSON.parse(message.data)
      switch (type) {
        case 'ListClients': {
          const newPeers = data.reduce((acc, c) => {
            if (c.identifier in this.peers) {
              acc[c.identifier] = this.peer[c.identifier]
            }
            return acc
          }, {})
          this.clients = data
          this.peers = newPeers
          break;
        }
        case 'SDP': {
          const { client: { identifier }, description } = data
          // if pc isnt in peers than its b/c its already waitng on response
          if (identifier in this.peers) {
            // can have peer but not pc if user declined to connect
            if ('pc' in this.peers[identifier]) {
              const { pc } = this.peers[identifier]
              pc.setRemoteDescription(description)
            }
          } else {
            this.peers[identifier] = { connectionRequest: description }
          }
          break;
        }
      }
    },
    //-----Socket-----
    newPC(identifier) {
      const pc = new PeerConnection(identifier)
      pc.descriptionCallback = this.sendDescription
      pc.fileCallback = this.downloadFile
      pc.stateCallback = this.pcStateChanged
      return pc
    },
    //-----PC-----
    connect ({ identifier }) {
      if (identifier in this.peers) {
        return // dont allow conneting to multiple
      }

      const pc = this.newPC(identifier)
      this.peers[identifier] = { pc }

      pc.createLocalDescription()
    },
    disconnect ({ identifier }) {
      if (identifier in this.peers) {
        const { pc } = this.peers[identifier]
        delete this.peers[identifier]
        pc.stop()
      }
    },
    sendDataMessage ({ identifier }) {
      const { pc } = this.peers[identifier] || {}
      if (!pc) {
        return
      }
      console.log('Sending file')
      pc.fileC.send(JSON.stringify({ hello: 'world' }))
    },
    //-----PC-----
    //-----PC callbacks-----
    sendDescription (identifier, description) {
      this.sendMessage('SDP', { identifier, description })
    },
    downloadFile(identifier, jsonFile) {
      console.log(jsonFile)
    },
    pcStateChanged (identifier, state) {
      console.log(`PC State: ${identifier} - ${state}`)
      if (identifier in this.peers) {
        switch(state) {
          case "new":
          case "checking":
          case "connecting":
            break;
          case "connected":
            this.peers[identifier].connected = true
            break;
          case "disconnected":
          case "closed":
            this.disconnect({ identifier })
            break;
          case "failed":
            this.peers[identifier].error = true
            break
        }
      }
    },
    //-----PC callbacks-----
    //-----Connected Request-----
    acceptConnection ({ identifier }) {
      const pc = this.newPC(identifier)
      const { connectionRequest } = this.peers[identifier]
      this.peers[identifier] = { pc }
      pc.setRemoteDescription(connectionRequest)
    },
    declineConnection ({ identifier }) {
      // dont delete identifier for peer so if the user keeps requesting
      delete this.peers[identifier].connectionRequest
    },
    //-----Connected Request-----
    sendMessage(type, data) {
      this.socket.send(JSON.stringify({ type, data }))
    }
  },
}
</script>

<style lang="scss">
#app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  align-items: center;
}
#app > * {
  max-width: 600px;
}
body {
  /* background-color: #2c3e50; */
  height: 100vh;
  margin: 0px;
}
</style>
