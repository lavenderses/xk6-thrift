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

  thrift.call(method, req);
}
