import request from "@/utils/request";

export function list(params) {
  return request({
    url: "appleDeveloper",
    method: "get",
    params,
  });
}

export function upload(data, config, cancel) {
  return request({
    url: "appleDeveloper",
    method: "post",
    data,
    onUploadProgress: config,
    cancelToken: cancel
  });
}

export function del(iss) {
  return request({
    url: "appleDeveloper?iss=" + iss,
    method: "delete"
  });
}

export function update(iss, limit, enable) {
  return request({
    url: "appleDeveloper?iss=" + iss + "&limit=" + limit + "&enable=" + enable,
    method: "patch"
  });
}
