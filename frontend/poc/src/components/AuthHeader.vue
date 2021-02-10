<template>
  <div class="navbar-menu">
    <div class="navbar-end">
      <div class="navbar-item">
        <div class="buttons">
          <a
            v-if="!account"
            @click="SignIn"
            target="_blank"
            rel="noopener noreferrer"
          >
            <i class="fas fa-sign-in-alt fa-2x" aria-hidden="false">Sign in</i>
          </a>
          <a v-else @click="SignOut" target="_blank" rel="noopener noreferrer">
            <i class="fas fa-sign-out-alt fa-2x" aria-hidden="false">Sign out</i>
          </a>
          <div v-if="account">{{ account.name }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { PublicClientApplication } from '@azure/msal-browser';
export default {
  name: 'HeaderBar',
  data() {
    return {
      account: undefined,
      github: 'https://github.com/cmatskas',
      signin: 'https://microsoft.com',
    };
  },
  async created() {
    this.$msalInstance = new PublicClientApplication(
      this.$store.state.msalConfig,
    );
  },
  mounted() {
    const accounts = this.$msalInstance.getAllAccounts();
    if (accounts.length == 0) {
      return;
    }
    this.account = accounts[0];
    this.$emit('login', this.account);

    this.$msalInstance.acquireTokenSilent(
      {
        scopes: ["api://24420be9-46e5-4584-acd7-64850d2f2a03/access_as_user"],
        account: this.account
      }).then(resp => {
        this.$emit('captoken', resp);
      }); 
  },
  methods: {
    async HandleMsalResponse(response) {
        console.log(response);
        const myAccounts = this.$msalInstance.getAllAccounts();
        this.account = myAccounts[0];
        this.$emit('login', this.account);
    },

    async SignIn() {

      await this.$msalInstance
        .loginPopup(this.$store.state.msalConfig.request)
        .then(() => {
          const myAccounts = this.$msalInstance.getAllAccounts();
          this.account = myAccounts[0];
          this.$emit('login', this.account);
          this.$msalInstance.acquireTokenSilent(
            {
              scopes: ["api://24420be9-46e5-4584-acd7-64850d2f2a03/access_as_user"],
              account: this.account
            }).then(resp => {
              this.$emit('captoken', resp);
            });
        })
        .catch(error => {
          console.error(`error during authentication: ${error}`);
        });
    },
    async SignOut() {
      await this.$msalInstance
        .logout({})
        .then(() => {
          this.$emit('logout');
        })
        .catch(error => {
          console.error(error);
        });
    },
  },
};
</script>

<style scoped>



a {
    height: 15px;
}

</style>