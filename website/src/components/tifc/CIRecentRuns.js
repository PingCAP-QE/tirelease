import * as React from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

export const CIRecentRuns =( {data, text} ) => {
    const [open, setOpen] = React.useState(false);
    const [scroll, setScroll] = React.useState('paper');
    const [maxWidth, setMaxWidth] = React.useState('lg');
    const descriptionElementRef = React.useRef(null);

    const handleClickOpen = (scrollType) => () => {
        setOpen(true);
        setScroll(scrollType);
    };
    const handleClose = () => {
        setOpen(false);
    };
    React.useEffect(() => {
        if (open) {
        const { current: descriptionElement } = descriptionElementRef;
        if (descriptionElement !== null) {
            descriptionElement.focus();
        }
        }
    }, [open]);

    return (
        <>
        <Button onClick={handleClickOpen('paper')} style={{justifyContent: "flex-start"}}>{text}</Button>
        {/* <Button onClick={handleClickOpen('body')}>scroll=body</Button> */}
        <Dialog
            open={open}
            onClose={handleClose}
            scroll={scroll}
            fullWidth={true}
            maxWidth={maxWidth}
            aria-labelledby="scroll-dialog-title"
            aria-describedby="scroll-dialog-description"
        >
            <DialogTitle id="scroll-dialog-title">Trace Logs</DialogTitle>
            <DialogContent dividers={scroll === 'paper'}>
                <DialogContentText
                    id="scroll-dialog-description"
                    ref={descriptionElementRef}
                    tabIndex={-1}
                >
                    <pre>{
                        JSON.stringify(data, null, 4) // Indented 4 spaces
                    }</pre>
                </DialogContentText>
            </DialogContent>
            <DialogActions>
            <Button variant="contained" onClick={handleClose}>Close</Button>
            </DialogActions>
        </Dialog>
        </>
    );
}
