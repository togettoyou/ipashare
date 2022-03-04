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

export function download(uuid) {
  return request({
    url: "download/ipa/" + uuid,
    method: "get",
    responseType: 'blob'
  });
}

export function del(uuid) {
  return request({
    url: "ipa?uuid=" + uuid,
    method: "delete"
  });
}
