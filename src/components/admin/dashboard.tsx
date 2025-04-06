import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { LayoutDashboard, Video, Users, Settings, Upload } from "lucide-react";
import UploadVideo from "./upload-video";

type Tab = "videos" | "upload" | "users" | "settings";

export default function AdminDashboard() {
  const [activeTab, setActiveTab] = useState<Tab>("videos");

  const renderContent = () => {
    switch (activeTab) {
      case "videos":
        return (
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Video Management</h2>
            <p>Manage your video content here.</p>
          </div>
        );
      case "upload":
        return <UploadVideo />;
      case "users":
        return (
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">User Management</h2>
            <p>Manage user accounts here.</p>
          </div>
        );
      case "settings":
        return (
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Settings</h2>
            <p>Configure system settings here.</p>
          </div>
        );
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="border-b">
        <div className="flex h-16 items-center px-4">
          <h1 className="text-2xl font-bold">Admin Dashboard</h1>
        </div>
      </div>
      <div className="flex">
        <div className="w-64 border-r p-4">
          <nav className="space-y-1">
            <Button
              variant={activeTab === "videos" ? "secondary" : "ghost"}
              className="w-full justify-start"
              onClick={() => setActiveTab("videos")}
            >
              <Video className="mr-2 h-4 w-4" />
              Videos
            </Button>
            <Button
              variant={activeTab === "upload" ? "secondary" : "ghost"}
              className="w-full justify-start"
              onClick={() => setActiveTab("upload")}
            >
              <Upload className="mr-2 h-4 w-4" />
              Upload Video
            </Button>
            <Button
              variant={activeTab === "users" ? "secondary" : "ghost"}
              className="w-full justify-start"
              onClick={() => setActiveTab("users")}
            >
              <Users className="mr-2 h-4 w-4" />
              Users
            </Button>
            <Button
              variant={activeTab === "settings" ? "secondary" : "ghost"}
              className="w-full justify-start"
              onClick={() => setActiveTab("settings")}
            >
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </Button>
          </nav>
        </div>
        <div className="flex-1 p-8">{renderContent()}</div>
      </div>
    </div>
  );
}
