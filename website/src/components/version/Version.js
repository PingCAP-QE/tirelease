import { useState } from "react";
import Stack from "@mui/material/Stack";

import { VersionTable } from "./VersionTable";
import { VersionAdd } from "./VersionAdd";
import { useQuery } from "react-query";
import { url } from "../../utils";
import { Button } from "@mui/material";

export default function Version() {
  const { isLoading, error, data } = useQuery("versions", () => {
    return fetch(url("version")).then(async (res) => {
      return await res.json();
    });
  });

  if (data) {
    console.log("version", data);
  }

  const [createVersion, setCreateVersion] = useState(false);

  return (
    <>
      {isLoading && <p>Loading...</p>}
      {error && <p>Error: {error.message}</p>}
      {data && (
        <Stack spacing={1}>
          <Stack direction={"row"} justifyContent={"flex-end"} spacing={2}>
            <Button
              variant="contained"
              onClick={() => {
                setCreateVersion(true);
              }}
            >
              Create version
            </Button>
          </Stack>
          <VersionTable data={data.data} />
        </Stack>
      )}
      {data && createVersion && (
        <VersionAdd
          open={createVersion}
          onClose={() => {
            setCreateVersion(false);
          }}
          versions={data.data.map((version) => version.name)}
        />
      )}
    </>
  );
}
