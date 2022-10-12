import request from "@/utils/request";

export function keyList(params) {
  return request({
    url: "key", method: "get",
    params
  });
}

export function addKey(data) {
  return request({
    url: "key", method: "post",
    data
  });
}

export function updateKey(data) {
  return request({
    url: "key/changenum", method: "post",
    data
  });
}

export function delKey(username) {
  return request({
    url: "key?username=" + username,
    method: "delete"
  });
}
