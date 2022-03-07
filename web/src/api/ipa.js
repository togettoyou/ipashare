import request from "@/utils/request";

export function list(params) {
  return request({
    url: "ipa",
    method: "get",
    params,
  });
}

export function upload(data, config, cancel) {
  return request({
    url: "ipa",
    method: "post",
    data,
    onUploadProgress: config,
    cancelToken: cancel
  });
}

export function download(uuid, config, cancel) {
  return request({
    url: "download/ipa/" + uuid,
    method: "get",
    responseType: 'blob',
    onDownloadProgress: config,
    cancelToken: cancel
  });
}

export function update(data) {
  return request({
    url: "ipa",
    method: "patch",
    data
  });
}

export function del(uuid) {
  return request({
    url: "ipa?uuid=" + uuid,
    method: "delete"
  });
}
