<template>
  <div id="app">
    <auth-header v-on:captoken="handleCapToken" v-on:login="login" v-on:logout="logout" />
    <HelloWorld v-bind:capabilities="this.capabilities" v-bind:costs="this.costs" v-bind:logged_in="this.logged_in"  msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import AuthHeader from './components/AuthHeader.vue'
import { mapMutations } from 'vuex';

export default {
    name: 'App',
    components: {
      HelloWorld,
      AuthHeader
    },
    data() {
      return {
        logged_in: false,
        capabilities: [],
        costs: {}
      }
    },
    mounted() {
      this.logged_in = false;
    },
    watch: {
      logged_in: function(val) {
        if (val === true) {
          console.log("x");
        }
      }
    },
    methods: {
    ...mapMutations(['setAccessToken']),
    login: function (account) {
        this.account = account;
        this.logged_in = true;
        console.log(this.account);
    },
    logout: function () {
        console.log('logging out');
        this.account = undefined;
        this.logged_in = false;
    },
    handleCapToken: function (resp) { 
      this.logged_in = true;
      this.getCapabilities(resp);
      this.getAllCosts(resp);
    },
    getCapabilities: function (token) {
      var req = new XMLHttpRequest();
      var x = this;
      req.onload = function () {
        var resp = JSON.parse(req.responseText);
        x.capabilities = resp.items;
        console.log(x.capabilities);
      };
      req.open("GET", "/api/get-capabilities");
      req.setRequestHeader("Authorization", "Bearer " + token.accessToken);
      req.send();

    },
    getAllCosts: function (token) {
      var req = new XMLHttpRequest();
      var x = this;
      req.onload = function () {
        var resp = JSON.parse(req.responseText);
        x.costs = resp;
        console.log(x.costs);

        var accountCostKv = {};
        resp.ResultsByTime[0].Groups.forEach(val => {
          accountCostKv[val.Keys[0].toString()] = val.Metrics.BlendedCost.Amount
        })
        x.costs = accountCostKv;
        console.log(x.costs);
      };
      req.open("GET", "/api/get-monthly-total-cost-all");
      req.setRequestHeader("Authorization", "Bearer " + token.accessToken);
      req.send();      
    }
  },
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  color: #2c3e50;
}
.navbar-menu {
  text-align: center;
}

html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
b, u, i, center,
dl, dt, dd, ol, ul, li,
fieldset, form, label, legend,
table, caption, tbody, tfoot, thead, tr, th, td,
article, aside, canvas, details, embed, 
figure, figcaption, footer, header, hgroup, 
menu, nav, output, ruby, section, summary,
time, mark, audio, video {
	margin: 0;
	padding: 0;
	border: 0;
	font-size: 100%;
	font: inherit;
	vertical-align: baseline;
}
/* HTML5 display-role reset for older browsers */
article, aside, details, figcaption, figure, 
footer, header, hgroup, menu, nav, section {
	display: block;
}
body {
	line-height: 1;
  background-color: #302b37;
}
ol, ul {
	list-style: none;
}
blockquote, q {
	quotes: none;
}
blockquote:before, blockquote:after,
q:before, q:after {
	content: '';
	content: none;
}
table {
	border-collapse: collapse;
	border-spacing: 0;
}
</style>
