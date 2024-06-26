// Copyright (c) 2022 Braden Nicholson


import axios from "axios";

const token = localStorage.getItem("token")

let headers = {
    headers: {
        Authorization: "Bearer " + token
    }
}


export default {
    async post(url: string, data?: {} | undefined): Promise<void> {
        const response = await axios.post(`https://api.udap.app${url}`, data, headers)
        let resp = response.data
        if (response.status !== 200) {
            // core.notify().show(`Request HTTPS ${response.status}`, resp, 2, 1000 * 8)
        }
    }
}

