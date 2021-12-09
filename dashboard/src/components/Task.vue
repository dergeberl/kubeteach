<template>
    <div id="info">
        <button @click="lastTask">back</button>

        Tasks:
        <select v-model="selectedTask" @change="renewStatus">
            <option v-for="task of tasks" :key="task.uid" :value="task.uid">
                {{ task.name }}
            </option>
        </select>

        <button @click="nextTask">next</button>

        <br>
        <div v-if="task">
            <h1>{{ task.title }}</h1>
            <p>
                {{ task.description }}
            </p>
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
    name: "Tasks",
    data() {
        return {
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