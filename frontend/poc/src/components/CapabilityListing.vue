<template>
  <div class="capabilityListing" v-if="this.logged_in">
      <div class="capability" v-for="capability in capabilities" :key="capability.id">
          <div style="display: flex">
              <div style="flex: 1">
                <span class="title entity">{{capability.name}}</span>
              </div>
              <div style="flex: 1; display: flex; justify-content: flex-end;">
                <span class="entity" style="margin-right: 10px;">${{ formatCost(costs[getAwsAccountId(capability)] )}}</span>
              </div>              
          </div>
          
      </div>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: String,
    logged_in: Boolean,
    capabilities: Array,
    costs: Object
  },
  data() {
      return {

      }
  },
  methods: {
      getAwsAccountId(cap) {
          if (cap.contexts.length > 0) {
              return cap.contexts[0].awsAccountId;
          } else {
              return 0;
          }
      },
      formatCost(val) {
          var x = parseFloat(val);
          if (Number.isNaN(x)) {
              return 0;
          }

          return x.toFixed(5);
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style  scoped>
h3 {
    margin: 40px 0 0;
}
ul {
    list-style-type: none;
    padding: 0;
    display: flex;
    flex-direction: column;
    font-size: 1.4em;
}
li {
    display: inline-block;
    margin: 0 10px;
    margin-top: 8px;
}
a {
    color: #42b983;
}

.title {
    margin-top: 4px;
    margin-bottom: 4px;
    margin-left: 10px;
}

.entity {
    font-size: 1.7em;
}

.capability {
    display: flex;
    flex-direction: column;
    font-size: 1.4em;
    padding-top: 8px;
    padding-bottom: 8px;
    margin-bottom: 8px;
    width: 100%;
    background-color: #7a7878;
}

.capabilityListing {
    margin-top: 10px;
    margin-left: 25px;
    margin-right: 25px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

</style>
