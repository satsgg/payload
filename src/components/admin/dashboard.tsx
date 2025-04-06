import React from 'react';
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { 
  LayoutDashboard, 
  Video, 
  Users, 
  Settings 
} from "lucide-react";

export default function AdminDashboard() {
  return (
    <div className="min-h-screen bg-background p-8">
      <h1 className="text-4xl font-bold mb-8">Admin Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card className="p-6">
          <div className="flex items-center gap-4">
            <Video className="h-8 w-8 text-primary" />
            <div>
              <h2 className="text-2xl font-semibold">Videos</h2>
              <p className="text-muted-foreground">Manage video content</p>
            </div>
          </div>
          <Button className="w-full mt-4">View Videos</Button>
        </Card>

        <Card className="p-6">
          <div className="flex items-center gap-4">
            <Users className="h-8 w-8 text-primary" />
            <div>
              <h2 className="text-2xl font-semibold">Users</h2>
              <p className="text-muted-foreground">Manage user accounts</p>
            </div>
          </div>
          <Button className="w-full mt-4">View Users</Button>
        </Card>

        <Card className="p-6">
          <div className="flex items-center gap-4">
            <LayoutDashboard className="h-8 w-8 text-primary" />
            <div>
              <h2 className="text-2xl font-semibold">Analytics</h2>
              <p className="text-muted-foreground">View site statistics</p>
            </div>
          </div>
          <Button className="w-full mt-4">View Analytics</Button>
        </Card>

        <Card className="p-6">
          <div className="flex items-center gap-4">
            <Settings className="h-8 w-8 text-primary" />
            <div>
              <h2 className="text-2xl font-semibold">Settings</h2>
              <p className="text-muted-foreground">Configure system</p>
            </div>
          </div>
          <Button className="w-full mt-4">View Settings</Button>
        </Card>
      </div>
    </div>
  );
}