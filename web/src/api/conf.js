import request from "@/utils/request";

export function getKeyConf() {
  return request({
    url: "conf/key", method: "get",
  });
}

export function setKeyConf(data) {
  return request({
    url: "conf/key", method: "post", data,
  });
}

export function get() {
  return request({
    url: "conf/oss", method: "get",
  });
}

export function set(data) {
  return request({
    url: "conf/oss", method: "post", data,
  });
}

export function verify() {
  return request({
    url: "conf/oss/verify", method: "get"
  });
}
