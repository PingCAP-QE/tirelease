import { url } from "../../../utils";

export function fetchVersion() {
  return fetch(url("version/maintained")).then(async (res) => {
    const data = await res.json();
    let { data: versions } = data;
    versions.sort();
    versions.reverse();
    return versions || [];
  });
}
