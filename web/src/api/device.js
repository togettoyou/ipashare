import request from "@/utils/request";

export function deviceList(iss) {
  return request({
    url: "appleDevice?iss=" + iss,
    method: "get"
  });
}

export function updateDevice(iss) {
  return request({
    url: "appleDevice?iss=" + iss,
    method: "post"
  });
}
