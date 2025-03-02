import { check, sleep } from 'k6';
import thrift from 'k6/x/thrift';
import ttypes from 'k6/x/thrift/ttypes';

export const options = {
  vus: 1,
  iterations: 1,
}

export default function() {
  const method = "simpleCall";
  const values = {};
  values[1] = ttypes.newTString("ID");
  const req = ttypes.newTRequest(values);

  const res = thrift.call(method, req);
  check(res, {
    "success Thrift call": (r) => r.isSuccess(),
  });

  sleep(0.5);

  const failValue = {};
  failValue[1] = ttypes.newTString("FAILURE");
  const failReq = ttypes.newTRequest(failValue);

  const failRes = thrift.call(method, failReq);
  check(failRes, {
    "failure Thrift call": (r) => !r.isSuccess(),
  });
}
