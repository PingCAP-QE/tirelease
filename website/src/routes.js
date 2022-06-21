import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import DataBoard from "./pages/databoard";
// import Triage from "./pages/triage";
import Assignments from "./pages/assignments";
import TiFC from "./pages/tifc";
import Release from "./pages/release";
import RecentOpen from "./pages/open";
import VersionPage from "./pages/version";
import AffectTriage from "./pages/affects";
import PickTriage from "./pages/close";
import AllIssues from "./pages/all";

const MyRoutes = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<VersionPage />} />
        <Route path="/home" element={<VersionPage />} />
        <Route path="/home/all" element={<AllIssues />} />
        <Route path="/home/open" element={<RecentOpen />} />
        <Route path="/home/affection" element={<AffectTriage />} />
        <Route path="/home/close" element={<PickTriage />} />
        <Route path="/home/databoard" element={<DataBoard />} />
        {/* <Route path="/home/triage" element={<Triage />} /> */}
        <Route path="/home/assignments" element={<Assignments />} />
        <Route path="/home/tifc" element={<TiFC />} />
        <Route path="/home/version" element={<VersionPage />} />
        <Route path="/home/triage" element={<Release />} />
        <Route path="/home/triage/:version" element={<Release />} />
      </Routes>
    </BrowserRouter>
  );
};

export default MyRoutes;
