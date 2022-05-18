import dayjs from "dayjs";

export function renderDateTime(time) {
    return (
      <>
        {time !== undefined && (
          <div>
            {dayjs(time).format(
              "YYYY-MM-DD HH:mm:ss"
            )}
          </div>
        )}
      </>
    );
}