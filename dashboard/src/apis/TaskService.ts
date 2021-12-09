import axios, {AxiosResponse} from "axios";

const API_ENDPOINT = "/api"
const TASK_ENDPOINT = API_ENDPOINT + "/tasks"
const TASK_STATUS_ENDPOINT = API_ENDPOINT + "/taskstatus"

export default new class TaskService {
    public listTasks() : Promise<Task[]> {
        return axios.get(TASK_ENDPOINT)
            .then(extractTasksFromAxiosResponse)
    }

    public getTaskStatus(taskID: string): Promise<TaskStatus> {
        return axios.get(TASK_STATUS_ENDPOINT + "/" + taskID)
            .then(extractTaskStatusFromAxiosResponse)
    }
}()

function extractTasksFromAxiosResponse(response: AxiosResponse): Task[] {
    return response.data
}

function extractTaskStatusFromAxiosResponse(response: AxiosResponse): TaskStatus {
    return response.data
}

export interface Task {
    name: string
    namespace: string
    title: string
    description: string
    uid: string
}

export interface TaskStatus {
    status: string
}
