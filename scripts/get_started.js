import thrift from 'k6/x/thrift';

export const options = {
  vus: 1,
  iterations: 3,
}

export default function() {
  thrift.echo();
}
