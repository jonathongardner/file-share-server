<template>
  <footer class="card-footer">
    <template v-if='pClient.connectionRequest'>
      <button class='button' @click='acceptConnection'>Accept</button>
      <button class='button' @click='declineConnection'>Decline</button>
    </template>
    <button v-else-if='pClient.error || pClient.connected' class="button" @click='disconnect'>Disconnect</button>
    <span v-else-if='pClient.connecting'  class="button">Connecting...</span>
    <button v-else class="button" @click='connect'>Connect</button>
  </footer>
</template>

<script>
export default {
  name: 'CardFooter',
  props: {
    pClient: {
      type: Object,
      required: true
    }
  },
  methods: {
    acceptConnection() {
      this.$emit('accept', this.pClient)
    },
    declineConnection() {
      this.$emit('decline', this.pClient)
    },
    disconnect() {
      this.$emit('disconnect', this.pClient)
    },
    connect() {
      this.$emit('connect', this.pClient)
    }
  }
}
</script>

<style scoped lang="scss">
</style>
