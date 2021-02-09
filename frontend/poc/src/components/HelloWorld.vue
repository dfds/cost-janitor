<template>
  <div class="container">
    <div class="inner">
    <h1>Billing</h1>
    <CapabilityListing />
    </div>
  </div>
</template>

<script>
import CapabilityListing from './CapabilityListing.vue'

export default {
  name: 'HelloWorld',
  props: {
    msg: String
  },
  components: {
    CapabilityListing
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

h1 {
  font-size: 4em;
}

.container {
    display: flex;
    flex-direction: column;
    margin: auto;
    width: 70%;
    align-items: flex-start;
    background-color: #423e48;
    color: #c9c9c9;
    height: 100vh;
}

.inner {
  margin-left: 35px;
  padding-top: 35px;
}
</style>
