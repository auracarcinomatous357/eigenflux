import { useList } from "@refinedev/core";
import { List } from "@refinedev/antd";
import { Table, Input, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useState } from "react";

interface Agent {
  agent_id: string;
  agent_name: string;
  email: string;
  bio: string;
  created_at: number;
  updated_at: number;
  profile_status: number | null;
  profile_keywords: string[];
}

const profileStatusMap: Record<number, { label: string; color: string }> = {
  0: { label: "Pending", color: "default" },
  1: { label: "Processing", color: "processing" },
  2: { label: "Failed", color: "error" },
  3: { label: "Completed", color: "success" },
};

const formatTimestamp = (ts: number) => {
  if (!ts) return "-";
  return new Date(ts).toLocaleString();
};

export const AgentList = () => {
  const [nameFilter, setNameFilter] = useState<string>("");
  const [current, setCurrent] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);

  const { query } = useList<Agent>({
    resource: "agents",
    pagination: {
      currentPage: current,
      pageSize,
      mode: "server",
    },
    filters: [...(nameFilter ? [{ field: "name", operator: "contains" as const, value: nameFilter }] : [])],
  });

  const columns: ColumnsType<Agent> = [
    {
      title: "ID",
      dataIndex: "agent_id",
      key: "agent_id",
      width: 80,
      fixed: "left",
    },
    {
      title: "Name",
      dataIndex: "agent_name",
      key: "agent_name",
      width: 150,
    },
    {
      title: "Email",
      dataIndex: "email",
      key: "email",
      width: 200,
    },
    {
      title: "Bio",
      dataIndex: "bio",
      key: "bio",
      width: 200,
      render: (text: string) => text ? (
        <Tooltip title={<div style={{ maxWidth: 400, whiteSpace: "pre-wrap", wordBreak: "break-all" }}>{text}</div>}>
          <Typography.Text
            copyable
            style={{ maxWidth: 200, display: "block" }}
            ellipsis
          >
            {text}
          </Typography.Text>
        </Tooltip>
      ) : "-",
    },
    {
      title: "Profile Status",
      dataIndex: "profile_status",
      key: "profile_status",
      width: 130,
      render: (status: number | null) => {
        if (status === null || status === undefined) return "-";
        const s = profileStatusMap[status];
        return s ? <Tag color={s.color}>{s.label}</Tag> : <Tag>{status}</Tag>;
      },
    },
    {
      title: "Profile Keywords",
      dataIndex: "profile_keywords",
      key: "profile_keywords",
      width: 200,
      render: (keywords: string[]) => {
        if (!keywords || keywords.length === 0) return "-";
        const joined = keywords.join(", ");
        return (
          <Tooltip title={<div style={{ maxWidth: 400, whiteSpace: "pre-wrap", wordBreak: "break-all" }}>{joined}</div>}>
            <Typography.Text
              copyable
              style={{ maxWidth: 200, display: "block" }}
              ellipsis
            >
              {joined}
            </Typography.Text>
          </Tooltip>
        );
      },
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      key: "created_at",
      width: 180,
      render: (ts: number) => formatTimestamp(ts),
    },
    {
      title: "Updated At",
      dataIndex: "updated_at",
      key: "updated_at",
      width: 180,
      render: (ts: number) => formatTimestamp(ts),
    },
  ];

  return (
    <List
      headerButtons={
        <>
          <Input.Search
            placeholder="Search by name"
            allowClear
            onSearch={(value) => {
              setNameFilter(value);
              setCurrent(1);
            }}
            style={{ width: 200, marginRight: 8 }}
          />
        </>
      }
    >
      <Table
        dataSource={query.data?.data}
        columns={columns}
        rowKey="agent_id"
        loading={query.isLoading}
        scroll={{ x: 1280 }}
        pagination={{
          current,
          pageSize,
          total: query.data?.total ?? 0,
          showSizeChanger: true,
          pageSizeOptions: [10, 20, 50, 100],
          onChange: (nextPage, nextPageSize) => {
            setCurrent(nextPage);
            setPageSize(nextPageSize);
          },
        }}
      />
    </List>
  );
};
