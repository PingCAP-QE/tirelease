const dev_config = {
  // The URL of the server to use.
  // SERVER_HOST: "http://localhost:8080/",
  SERVER_HOST: "http://tirelease.pingcap.net/",
};

const prod_config = {
  // The URL of the server to use.
  SERVER_HOST: "/",
};

console.log(process.env);

export default process.env.NODE_ENV === "development"
  ? dev_config
  : prod_config;

export const GIT_CLIENT_ID = "(GIT_CLIENT_ID)";
export const GIT_CLIENT_SECRET = "(GIT_CLIENT_SECRET)";
