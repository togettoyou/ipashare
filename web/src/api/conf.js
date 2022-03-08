import request from "@/utils/request";

export function get() {
  return request({
    url: "conf/oss",
    method: "get",
  });
}

export function set(data) {
  return request({
    url: "conf/oss",
    method: "post",
    data,
  });
}

export function verify() {
  return request({
    url: "conf/oss/verify",
    method: "get"
  });
}
