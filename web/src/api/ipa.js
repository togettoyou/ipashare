import request from "@/utils/request";

export function list(params) {
  return request({
    url: "ipa",
    method: "get",
    params,
  });
}

export function upload(data) {
  return request({
    url: "ipa",
    method: "post",
    data,
  });
}

export function del(uuid) {
  return request({
    url: "ipa?uuid=" + uuid,
    method: "delete"
  });
}
