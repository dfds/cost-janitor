<template>
  <div class="container">
    <div class="inner">
      <h1>Billing</h1>
      <CapabilityListing v-bind:capabilities="this.capabilities" v-bind:costs="this.costs" v-bind:logged_in="this.logged_in" />
    </div>
  </div>
</template>

<script>
import CapabilityListing from './CapabilityListing.vue'

export default {
  name: 'HelloWorld',
  props: {
    msg: String,
    logged_in: Boolean,
    capabilities: Array,
    costs: {}
  },
  data() {
    return {
    }
  },
  components: {
    CapabilityListing
  }
}


</script>

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
  margin-left: 15px;
}

.container {
    display: flex;
    flex-direction: column;
    margin: auto;
    width: 70%;
    align-items: flex-start;
    background-color: #423e48;
    color: #c9c9c9;
    padding-bottom: 15px;
}

.inner {
    padding-top: 35px;
    display: flex;
    flex-direction: column;
    width: 100%;
}
</style>
