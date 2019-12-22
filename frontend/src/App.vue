<template>
  <div id="app">
    <div v-if="joinedLobby !== null && joinedLobby.InProgress">
      <h1>You are <span :style="'color: ' + colors[playerID]">Player {{playerID}}</span> in game {{joinedLobby.LobbyID}} </h1>
      <p>Category: {{gameData.category}}</p>
      <p>
        <span v-if="gameData.fake === playerID">You are the Fake</span>
        <span v-else>Word: {{gameData.name}}</span>
      </p>
      <button v-on:click="exitLobby">Exit Lobby</button>
      <button v-if="playerID === 1" v-on:click="updateLobbyData">New Word</button>
    </div>
    <div v-else-if="joinedLobby !== null">
      <h1>You are in game {{joinedLobby.LobbyID}} </h1>
      <ul>
        <li :style="'color: ' + colors[player]" v-for="(player, i) in joinedLobby.Players" :key="i">Player {{player}} is {{colors[player]}} <span v-if="player === playerID" style="color: black">← This is you</span></li>
      </ul>
      <button v-on:click="exitLobby">Exit Lobby</button>
      <button v-if="playerID === 1" v-on:click="startGame">Start</button>
    </div>
    <div v-else>
      <h1>Open Lobbies</h1>
      <ul>
        <li v-for="(lobby, i) in openLobbie" :key="i">Lobby number {{lobby.LobbyID}} <button v-on:click="joinLobby(lobby.Hashid)">Join</button></li>
      </ul>
      <button v-on:click="createLobby">Create Lobby</button>
    </div>
  </div>
</template>

<script>
/* eslint-disable */
import axios from 'axios';
import {wordsEN} from './words-en.js';

export default {
  name: 'app',
  data() {
    return {
      colors: ['', 'black', 'red', 'green', 'blue', 'orange', 'yellow', 'darkblue', 'brown', 'darkred', 'darkgreen'],
      lobbies: [],
      joinedLobby: null,
      playerID: 0,
      pollingLobby: null,
      gameData: {category: '', word: '', fake: 0},
    }
  },
  mounted() {
    this.listLobbies();
    var s = window.localStorage.getItem('save');
    if (s !== null) {
      var save = JSON.parse(s);
      this.playerID = save.playerID;
      this.getLobbyData(save.hashid);
      this.startPollingLobby();
    }
  },
  beforeDestroy() {
    this.stopPollingLobby();
  },
  computed: {
    openLobbie() {
      return this.lobbies.filter((l) => {
        return l.InProgress === false;
      });
    }
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
    getLobbyData(hashid) {
      axios.get('/lobby/' + hashid).then((response) => {
        this.joinedLobby = response.data;
        if (this.joinedLobby.Data !== '') {
          this.gameData = JSON.parse(this.joinedLobby.Data);
        }
      });
    },
    joinLobby(hashid) {
      axios.get('/join/' + hashid).then((response) => {
        this.joinedLobby = response.data;
        this.playerID = this.joinedLobby.YourPlayerID;
        var s = {playerID: this.joinedLobby.YourPlayerID, lobbyID: this.joinedLobby.LobbyID, hashid: this.joinedLobby.Hashid};
        window.localStorage.setItem('save', JSON.stringify(s));
        this.startPollingLobby();
      });
    },
    exitLobby() {
      var ok = confirm('Are you sure?');
      if (ok) {
        window.localStorage.removeItem('save');
        location.reload();
      }
    },
    startPollingLobby(hashid) {
      this.pollingLobby = setInterval(() => {
        this.getLobbyData(this.joinedLobby.Hashid);
      }, 5000);
    },
    stopPollingLobby() {
      clearInterval(this.pollingLobby);
      this.pollingLobby = null;
    },
    startGame() {
      this.updateLobbyData();
    },
    updateLobbyData() {
      // get all things to play, filter out played, send data
      var pos = randomIntFromInterval(0, wordsEN.length-1);
      var card = wordsEN[pos];

      axios.post('/update/' + this.joinedLobby.Hashid, {
        category: card.category,
        name: card.name,
        fake: randomIntFromInterval(1, this.joinedLobby.Players.length),
      })
      .then(function (response) {
        console.log(response);
      })
    }
  }
}

function randomIntFromInterval(min, max) {
	return Math.floor(Math.random()*(max-min+1)+min);
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

h1 {
  font-size: 24px;
}
p {
  font-size: 20px;
}

ul {
  text-align: left;
  width: 300px;
  margin: auto;
  margin-bottom: 20px;
}
</style>
