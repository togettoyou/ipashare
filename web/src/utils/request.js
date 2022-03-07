import axios from "axios";
import {Message} from "element-ui";
import store from "@/store";
import {getToken} from "@/utils/auth";

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 0, // request timeout
});

// request interceptor
service.interceptors.request.use(
  (config) => {
    if (store.getters.token) {
      config.headers["Authorization"] = getToken();
    }
    return config;
  },
  (error) => {
    console.log(error);
    return Promise.reject(error);
  }
);

// response interceptor
service.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code) {
      if (res.code !== 0) {
        Message({
          message: res.msg || "Error",
          type: "error",
          duration: 5 * 1000,
        });
        return Promise.reject(new Error(res.message || "Error"));
      }
    }
    return res;
  },
  (error) => {
    const res = error.response.data
    // 令牌失效,或令牌不存在
    if (res.code === 20104 || res.code === 20105 || res.code === 20106) {
      location.href = "/admin/#/login";
    }
    Message({
      message: res.msg,
      type: "error",
      duration: 5 * 1000,
    });
    return Promise.reject(error);
  }
);

export default service;
