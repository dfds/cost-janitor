import Vuex from 'vuex';
import Vue from 'vue'

Vue.use(Vuex)

const store = new Vuex.Store({
    state: {
        msalConfig: {
          auth: {
            clientId: '41b349a5-8c08-4c98-a4f1-e0896d320a9e',
            authority: 'https://login.microsoftonline.com/73a99466-ad05-4221-9f90-e7142aa2f6c1',
            redirectUri: "https://cost-janitor.hellman.oxygen.dfds.cloud"
          },
          cache: {
            cacheLocation: 'localStorage',
          },
          request: {
            scopes: ["user.read", "offline_access", "openid", "api://24420be9-46e5-4584-acd7-64850d2f2a03/access_as_user"],
            redirectUri: "https://cost-janitor.hellman.oxygen.dfds.cloud"
        }
        },
        accessToken: ''
    },
    mutations: {
        setAccessToken(state, token){
          state.accessToken = token;
        } 
    }
});

export default store;