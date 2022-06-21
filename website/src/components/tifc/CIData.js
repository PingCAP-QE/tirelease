import { useEffect, useState } from "react";
import Stack from "@mui/material/Stack";
import Button from '@mui/material/Button';

import { CIJobNameSelector} from "./CIJobNameSelector";
import { CIDatePicker} from "./CIDatePicker";
import { RenderCITable} from "./CITable";
import { CITips } from "./CITips";
import { RenderCITableGrid} from "./CITableGrid";

function addDays(date, days) {
    var result = new Date(date);
    result.setDate(result.getDate() + days);
    return result;
}

export default function CIData() {
    const [jobName, setJobName] = useState("tidb-unit-test-hourly");
    const [timestamp, setTimestamp] = useState(addDays(new Date(), -3));
    const [table, setTable] = useState(null);
    const refreshTable = () => {
        fetch("http://172.16.5.15:30792/report/unstable?timestamp=" + Math.round(timestamp.getTime()/1000))
        .then(response => response.json())
        .then(data => {
            setTable(data);
        })
        .catch((e) => {
            console.log(e);
        });
    }
    useEffect(() => {
        refreshTable();
    }, []); // empty dependency array

    return (
        <>
        <Stack spacing={1}>
            <Stack direction={"row"} justifyContent={"flex-start"} spacing={2}>
                {/* <CIJobNameSelector jobName={jobName} setJobName={setJobName}/> */}
                <CIDatePicker timestamp={timestamp} setTimestamp={setTimestamp}/>
                <Button variant="contained" onClick={refreshTable}>Query</Button>
                <CITips />
            </Stack>
            {table && (
                <Stack direction={"row"} justifyContent={"flex-start"} spacing={2}>
                    {/* <RenderCITable data={table} /> */}
                    <RenderCITableGrid data={table}/>
                </Stack>
            )}
        </Stack>
        </>
    );
}
