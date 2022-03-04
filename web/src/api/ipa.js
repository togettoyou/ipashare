import request from "@/utils/request";

export function list(params) {
  return request({
    url: "ipa",
    method: "get",
    params,
  });
}

export function upload(data, config) {
  return request({
    url: "ipa",
    method: "post",
    data,
    onUploadProgress: config
  });
}

export function download(uuid, config) {
  return request({
    url: "download/ipa/" + uuid,
    method: "get",
    responseType: 'blob',
    onDownloadProgress: config
  });
}

export function del(uuid) {
  return request({
    url: "ipa?uuid=" + uuid,
    method: "delete"
  });
}
