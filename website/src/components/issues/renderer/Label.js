import { Chip } from "@mui/material";
export function getLabelValue(filter, mapper) {
  return (params) => {
    const filtered = params.row.issue.labels.filter(filter);
    return filtered
      .map((label) => label.name)
      .map(mapper)
      .join(",");
  };
}

export function renderLabel(filter, mapper) {
  return (params) => {
    const filtered = params.row.issue.labels.filter(filter);
    return filtered.map((label) => (
      //   <Chip label={mapper(label.name)} color={"#" + label.color}></Chip>
      <Chip
        label={mapper(label.name)}
        size="small"
        style={{ backgroundColor: "#" + label.color }}
      ></Chip>
    ));
  };
}
