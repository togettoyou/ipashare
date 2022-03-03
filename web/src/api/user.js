import request from "@/utils/request";

export function login(data) {
  return request({
    url: "user/login",
    method: "post",
    data,
  });
}

export function changePW(data) {
  return request({
    url: "user/changepw",
    method: "post",
    data,
  });
}
