import React, { useState } from 'react';
import { Form, Input, Button, Tabs, message } from 'antd';
import type { TabsProps } from 'antd';
import type { FormInstance } from 'antd/es/form';
import type { Rule } from 'antd/es/form';
import { userApi } from '../services/api';
import { LoginResponse, SignupResponse } from '../types';

interface FormValues {
  type: 'login' | 'register';
  username: string;
  password: string;
  confirmPassword?: string;
}

const Auth: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm<FormValues>();

  const handleSubmit = async (values: FormValues) => {
    try {
      setLoading(true);
      const { username, password } = values;
      const isLogin = form.getFieldValue('type') === 'login';

      if (isLogin) {
        const { data: responseData } = await userApi.login(username, password);
        if (responseData.code === 0 && responseData.data.token) {
          localStorage.setItem('token', responseData.data.token);
          message.success('登录成功！');
          window.location.reload();
        } else {
          message.error(responseData.message || '登录失败');
        }
      } else {
        const { data: responseData } = await userApi.signup(username, password);
        if (responseData.code === 0) {
          message.success('注册成功！请登录');
          form.setFieldsValue({ type: 'login' });
        } else {
          message.error(responseData.message || '注册失败');
        }
      }
    } catch (error) {
      console.error('操作失败:', error);
      message.error('操作失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const validateConfirmPassword: Rule = ({ getFieldValue }) => ({
    validator(_, value) {
      if (!value || getFieldValue('password') === value) {
        return Promise.resolve();
      }
      return Promise.reject(new Error('两次输入的密码不一致'));
    },
  });

  const items: TabsProps['items'] = [
    {
      key: 'login',
      label: '登录',
      children: (
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{ type: 'login' }}
        >
          <Form.Item name="type" hidden>
            <Input />
          </Form.Item>

          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              登录
            </Button>
          </Form.Item>
        </Form>
      ),
    },
    {
      key: 'register',
      label: '注册',
      children: (
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{ type: 'register' }}
        >
          <Form.Item name="type" hidden>
            <Input />
          </Form.Item>

          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认密码"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认密码' },
              validateConfirmPassword,
            ]}
          >
            <Input.Password placeholder="请确认密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              注册
            </Button>
          </Form.Item>
        </Form>
      ),
    },
  ];

  return (
    <div style={{ maxWidth: 400, margin: '0 auto', padding: '24px' }}>
      <Tabs defaultActiveKey="login" items={items} />
    </div>
  );
};

export default Auth; 