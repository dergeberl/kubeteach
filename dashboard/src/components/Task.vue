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

<script lang="ts">
import {defineComponent} from 'vue'
import TaskService, {Task, TaskStatus} from "@/apis/TaskService"

export default defineComponent({
    name: "Tasks",
    data() {
        return {
            tasks: [] as Task[],
            selectedTask: "" as string,
            selectedTaskStatus: "" as string,
            interval: null as any
        };
    },
    computed: {
        task: function () : Task {
            if (typeof this.tasks !== "object")
                return {} as Task

            let tasks = this.tasks.filter((task: Task) => task.uid === this.selectedTask)

            if (!tasks) {
                return {} as Task
            }

            return tasks[0]
        }
    },
    mounted() {
        TaskService.listTasks()
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
        cancelFetchTaskStatusInterval() : void {
            if (this.interval) {
                clearTimeout(this.interval)
            }
        },
        setFetchTaskStatusInterval() : void {
            // Cancel the old interval, to prevent multiple concurrent intervals
            if (this.interval) {
                this.cancelFetchTaskStatusInterval()
            }
            this.interval = setInterval(this.getStatus, 2000)
        },
        saveTasksToData(tasks: Task[]) : Task[] {
            this.tasks = tasks
            return tasks
        },
        cleanStatus() {
            return new Promise((resolve => {
                this.selectedTaskStatus = ""
                resolve(null)
            }))
        },
        getStatus() {
            if (this.selectedTask) {
                return TaskService.getTaskStatus(this.selectedTask)
                    .then((taskStatus : TaskStatus) => this.selectedTaskStatus = taskStatus.status)
                    .catch(e => console.error(e))
            }
            return new Promise(((resolve) => resolve(null)))
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
        selectFirstTaskIfNoneSelected: function (tasks: Task[]) {
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
})

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