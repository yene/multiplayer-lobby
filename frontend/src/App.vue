<template>
  <div id="app">
    <div v-if="joinedLobby !== null">
      <h1>You are in game {{joinedLobby.LobbyID}} </h1>
      <ul>
        <li :style="'color: ' + colors[player]" v-for="(player, i) in joinedLobby.Players" :key="i">Player {{colors[player]}} </li>
      </ul>
      <button v-on:click="createLobby">Create Lobby</button>
    </div>
    <div v-else>
      <h1>Games in progress</h1>
      <ul>
        <li v-for="(lobby, i) in lobbies" :key="i">Lobby number {{lobby.LobbyID}} <button v-on:click="joinLobby(lobby.Hashid)">Join</button></li>
      </ul>
      <button v-on:click="createLobby">Create Lobby</button>
    </div>
  </div>
</template>

<script>
/* eslint-disable */
import axios from 'axios';

export default {
  name: 'app',
  components: {
  },
  data() {
    return {
      colors: ['black', 'darkred', 'red', 'orange', 'yellow', 'green', 'darkgreen', 'blue', 'darkblue', 'brown'],
      lobbies: [],
      joinedLobby: null,
    }
  },
  mounted() {
    this.listLobbies();
  },
  methods: {
    createLobby() {
      axios.get('/create').then((response) => {
        this.listLobbies();
        this.joinLobby(response.data.Hashid);
      });
    },
    listLobbies() {
      axios.get('/lobbies').then((response) => {
        var l = Object.values(response.data);
        l.sort((a,b) => {
          return a.LobbyID - b.LobbyID;
        })
        this.lobbies = l;
      });
    },
    joinLobby(hashid) {
      axios.get('/join/' + hashid).then((response) => {
        this.joinedLobby = response.data;
        var s = {playerid: this.joinedLobby.YourPlayerID, lobbyid: this.joinedLobby.LobbyID};
        window.localStorage.setItem('gamestate', JSON.stringify(s));
      });
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 20px;
}
</style>
