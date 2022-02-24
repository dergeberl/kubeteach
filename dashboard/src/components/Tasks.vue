<template>
  <v-app id="kubeteach">

    <v-app-bar app flat>
    <v-app-bar-nav-icon @click="showMenue = !showMenue"></v-app-bar-nav-icon>
    <v-toolbar-title>KUBETEACH</v-toolbar-title>
    <v-spacer></v-spacer>
      <v-btn icon>
        <a href="https://www.github.com/dergeberl/kubeteach" target="_blank" rel="noreferrer noopener" style="all: unset;">
          <v-icon>mdi-github</v-icon>
        </a>
      </v-btn>
    </v-app-bar>

    <v-navigation-drawer v-if="showMenue" app>
     <v-list-item  @click="selectedTask = task.uid; getStatus()" v-for="task of tasks" :key="task.uid" :value="task.uid" link>
        <v-list-item-content>
          <v-list-item-title>{{ task.name }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-navigation-drawer>
    <v-main>
      <v-card
        class="d-flex justify-space-between pa-2"
        height=100%>
        <div v-if="task" style="height: 100%; width: 40%; text-align: center" class="">
          <h1>{{ task.title }}</h1>
          <p>
              {{ task.description }}
          </p>
          Status: {{ selectedTaskStatus }}
        </div> 
        <div style="height: 100%; width: 60%">
          <iframe src="/shell" style="height: 100%; width:100%; borders: 0" />
        </div>
      </v-card>
    </v-main>

  </v-app>
</template>

<script>
import axios from "axios";

let apiUrl = "/api/"

function extractResponseFromAxios(response) {
    return response.data
}

function fetchTaskStatus(taskID) {
    return axios.get(apiUrl + `taskstatus/` + taskID)
        .then(extractResponseFromAxios)
}

function fetchTasks() {
    return axios.get(apiUrl + `tasks`)
        .then(extractResponseFromAxios)
}

export default {
    name: "KubeteachTasks",
    data() {
        return {
            showMenue: true,
            tasks: [],
            selectedTask: "",
            selectedTaskStatus: "",
            interval: null
        };
    },
    computed: {
        task: function () {
            if (typeof this.tasks !== "object")
                return null

            let tasks = this.tasks.filter(task => task.uid === this.selectedTask)

            if (!tasks) {
                return null
            }

            return tasks[0]
        }
    },
    mounted() {
        fetchTasks()
            .then(this.saveTasksToData)
            .then(this.selectFirstTaskIfNoneSelected)
            .catch(e => console.error(e))
            .then(this.setFetchTaskStatusInterval)
    },
    unmounted() {
        this.cancelFetchTaskStatusInterval()
    },
    methods: {
        renewStatus() {
            return this.cleanStatus().then(this.getStatus)
        },
        cancelFetchTaskStatusInterval() {
            if (this.interval) {
                clearTimeout(this.interval)
            }
        },
        setFetchTaskStatusInterval() {
            // Cancel the old interval, to prevent multiple concurrent intervals
            if (this.interval) {
                this.cancelFetchTaskStatusInterval()
            }
            this.interval = setInterval(this.getStatus, 2000)
        },
        saveTasksToData(tasks) {
            this.tasks = tasks
            return tasks
        },
        cleanStatus() {
            return new Promise((resolve => {
                this.tasks.selectedTaskStatus = ""
                resolve()
            }))
        },
        getStatus() {
            if (this.selectedTask) {
                return fetchTaskStatus(this.selectedTask)
                    .then(taskStatus => this.selectedTaskStatus = taskStatus.status)
                    .catch(e => console.error(e))
            }
            return new Promise(((resolve) => resolve()))
        },
        nextTask() {
            let found = false
            this.tasks.forEach(t => {
                if (found) {
                    this.selectedTask = t.uid
                    this.renewStatus()
                    found = false
                } else if (this.selectedTask === t.uid) {
                    found = true
                }
            });
        },
        selectFirstTaskIfNoneSelected: function (tasks) {
            if (!this.selectedTask) {
                this.selectedTask = tasks[0].uid
                this.getStatus()
            }
        },
        lastTask() {
            let last = ""
            this.tasks.forEach(t => {
                if (this.selectedTask === t.uid && last !== "") {
                    this.selectedTask = last
                    this.renewStatus()
                }
                last = t.uid
            });
        }
    }
}

</script>

<style scoped>

</style>