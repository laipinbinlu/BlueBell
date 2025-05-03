import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Select, message, Card, Space } from 'antd';
import { communityApi, postApi } from '../services/api';
import { Community } from '../types';
import { handleApiError, handleApiResponse } from '../utils/errorHandler';
import { formStyles, buttonStyles } from '../styles/common';

interface CreatePostProps {
  onSuccess?: () => void;
}

interface PostFormData {
  title: string;
  content: string;
  community_id: number;
}

const CreatePost: React.FC<CreatePostProps> = ({ onSuccess }) => {
  const [form] = Form.useForm<PostFormData>();
  const [communities, setCommunities] = useState<Community[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetchCommunities();
  }, []);

  const fetchCommunities = async () => {
    try {
      setLoading(true);
      const response = await communityApi.getCommunityList();
      handleApiResponse(
        response,
        (data) => setCommunities(data),
        '获取社区列表失败'
      );
    } catch (error) {
      handleApiError(error, '获取社区列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (values: PostFormData) => {
    try {
      setSubmitting(true);
      const response = await postApi.createPost(values);
      
      if (handleApiResponse(response, () => {
        message.success('发帖成功！');
        form.resetFields();
        onSuccess?.();
      }, '发帖失败')) {
        // 额外的成功处理逻辑（如果需要）
      }
    } catch (error) {
      handleApiError(error, '发帖失败，请重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Card title="发布新帖子" style={formStyles}>
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        requiredMark="optional"
      >
        <Form.Item
          name="title"
          label="标题"
          rules={[
            { required: true, message: '请输入标题' },
            { max: 100, message: '标题最多100个字符' }
          ]}
        >
          <Input placeholder="请输入帖子标题" />
        </Form.Item>

        <Form.Item
          name="community_id"
          label="社区"
          rules={[{ required: true, message: '请选择社区' }]}
        >
          <Select
            placeholder="请选择社区"
            loading={loading}
            disabled={loading}
            showSearch
            optionFilterProp="children"
          >
            {communities.map((community) => (
              <Select.Option key={community.id} value={community.id}>
                {community.name}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="content"
          label="内容"
          rules={[
            { required: true, message: '请输入内容' },
            { max: 5000, message: '内容最多5000个字符' }
          ]}
        >
          <Input.TextArea
            rows={6}
            placeholder="请输入帖子内容"
            showCount
            maxLength={5000}
          />
        </Form.Item>

        <Form.Item>
          <Space>
            <Button
              type="primary"
              htmlType="submit"
              loading={submitting}
              style={buttonStyles}
            >
              发布帖子
            </Button>
            <Button onClick={() => form.resetFields()}>
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default CreatePost; 