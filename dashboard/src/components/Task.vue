<template>
  <div id="info">
    <button @click="lastTask()">back</button> 

    Tasks: 
    <select v-model="selectedTask" @change="cleanStatus">
      <option v-for="task of tasks" :key="task.uid" :value="task.uid">
        {{ task.name }}
      </option>
    </select>

    <button @click="nextTask()">next</button>

    <br>
    <div v-for="task of tasks" :key="task.uid">
      <div v-if="selectedTask==task.uid">
        <h1>{{ task.title }}</h1>
        <p>
          {{ task.description }}
        </p>
      </div>
    </div>
    <br>
    Status: {{ selectedTaskStatus }}
    <br>
  </div>
  <iframe src="/shell" id="shell"></iframe>
</template>

<script>
import axios from "axios";

let apiUrl = "/api/"
export default {
  name: "Tasks",
  data() {
    return {
      tasks: [],
      selectedTask: "",
      selectedTaskStatus: ""
    };
  },
  async created() {
    try {
      const res = await axios.get(apiUrl+`tasks`);
      this.tasks = res.data;
      if (this.selectedTask == "") {
        this.selectedTask = this.tasks[0].uid
        this.getStatus()
      }
    } catch (e) {
      console.error(e);
    }
    this.timer = setInterval(this.getStatus, 2000);
  },
  methods: {
    cleanStatus() {
      this.tasks.selectedTaskStatus= ""
      this.getStatus()
    },
    async getStatus() {
      if (this.selectedTask != "") {
        try {
          const res = await axios.get(apiUrl+`taskstatus/`+this.selectedTask);
          this.selectedTaskStatus = res.data.status;
        } catch (e) {
          console.error(e);
        }
      } 
    },
    nextTask() {
      let found = false
      this.tasks.forEach(t => {
        if (found) {
          this.selectedTask = t.uid
          this.cleanStatus()
          found = false
        }else if (this.selectedTask == t.uid) {
          found = true
        }
      });
    },
    lastTask() {
      let last = ""
      this.tasks.forEach(t => {
        if (this.selectedTask == t.uid && last != "") {
          this.selectedTask = last
          this.cleanStatus()
        }
        last = t.uid
      });
    }
  }
}
  
</script>

<style scoped>
#info {
  margin: 10px;
  position: absolute;
  top: 0px;
  left: 0px;
  width: 40%;
  height: 90%;
  padding: 10px;
}
#shell {
  border: none;
  position: absolute;
  top: 0px;
  right: 0px;
  width: 60%;
  height: 100%;
  text-align: left;
}
</style>