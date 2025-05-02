import React, { useEffect, useState } from 'react';
import { Form, Input, Button, Select, message, Space } from 'antd';
import { Community } from '../types';
import { postApi, communityApi } from '../services/api';

const { TextArea } = Input;

interface PostFormData {
  title: string;
  content: string;
  community_id: number;
}

const CreatePost: React.FC<{ onSuccess?: () => void }> = ({ onSuccess }) => {
  const [form] = Form.useForm<PostFormData>();
  const [communities, setCommunities] = useState<Community[]>([]);
  const [loading, setLoading] = useState(false);
  const [fetchingCommunities, setFetchingCommunities] = useState(false);

  useEffect(() => {
    const fetchCommunities = async () => {
      try {
        setFetchingCommunities(true);
        const response = await communityApi.getCommunityList();
        if (response.code === 0) {
          setCommunities(response.data);
        } else {
          message.error(response.message || '获取社区列表失败');
        }
      } catch (error) {
        console.error('获取社区列表失败:', error);
        message.error('获取社区列表失败，请稍后重试');
      } finally {
        setFetchingCommunities(false);
      }
    };
    fetchCommunities();
  }, []);

  const handleSubmit = async (values: PostFormData) => {
    try {
      setLoading(true);
      const response = await postApi.createPost(values);
      if (response.code === 0) {
        message.success('发帖成功！');
        form.resetFields();
        onSuccess?.();
      } else {
        message.error(response.message || '发帖失败');
      }
    } catch (error) {
      console.error('发帖失败:', error);
      message.error('发帖失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={handleSubmit}
      style={{ maxWidth: 800, margin: '0 auto', padding: '24px' }}
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
        label="选择社区"
        rules={[{ required: true, message: '请选择社区' }]}
      >
        <Select
          placeholder="请选择发帖社区"
          loading={fetchingCommunities}
          disabled={fetchingCommunities}
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
        <TextArea
          rows={6}
          placeholder="请输入帖子内容"
          showCount
          maxLength={5000}
        />
      </Form.Item>

      <Form.Item>
        <Space>
          <Button type="primary" htmlType="submit" loading={loading}>
            发布帖子
          </Button>
          <Button onClick={() => form.resetFields()}>重置</Button>
        </Space>
      </Form.Item>
    </Form>
  );
};

export default CreatePost; 