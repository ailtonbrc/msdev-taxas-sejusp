import { Select, Checkbox, Button, Divider } from 'antd';
import type { SelectProps } from 'antd';
import { ClearOutlined } from '@ant-design/icons';

interface SelectTableauOption {
  label: string;
  value: string;
}

interface SelectTableauProps extends Omit<SelectProps, 'options'> {
  options: (string | SelectTableauOption)[];
  loading?: boolean;
}

const SelectTableau: React.FC<SelectTableauProps> = ({ options, value, onChange, placeholder, ...rest }) => {
  const isMultiple = rest.mode === 'multiple';


  const normalizedOptions: SelectTableauOption[] = options.map(opt =>
    typeof opt === 'string' ? { label: opt, value: opt } : opt
  );

  const handleLimparLocal = () => {
    if (onChange) {
      onChange(isMultiple ? [] : undefined, []);
    }
  };

  const temValor = Array.isArray(value) ? value.length > 0 : !!value;

  return (
    <Select
      {...rest}
      value={value}
      onChange={onChange}
      placeholder={placeholder}
      maxTagCount={0}
      maxTagPlaceholder={(omittedValues) => {
        if (omittedValues.length === 1) {
          return omittedValues[0].label;
        }
        return `(Múltiplos - ${omittedValues.length})`;
      }}
      optionLabelProp="label"
      dropdownRender={(menu) => (
        <>
          <div style={{ padding: '4px 8px', display: 'flex', justifyContent: 'flex-start' }}>
            <Button
              type="link"
              size="small"
              icon={<ClearOutlined />}
              onClick={handleLimparLocal}
              disabled={!temValor}
            >
              Limpar
            </Button>
          </div>
          <Divider style={{ margin: '2px 0' }} />
          {menu}
        </>
      )}
    >
      {normalizedOptions.map(opt => {
        const selected = Array.isArray(value) ? value.includes(opt.value) : value === opt.value;

        return (
          <Select.Option key={opt.value} value={opt.value} label={opt.label}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '2px 0' }}>
              {isMultiple && (
                <Checkbox
                  checked={selected}
                  style={{ pointerEvents: 'none' }}
                />
              )}
              <span style={{
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                color: selected ? '#1890ff' : 'inherit',
                fontWeight: selected ? 500 : 400
              }}>
                {opt.label}
              </span>
            </div>
          </Select.Option>
        );
      })}
    </Select>
  );
};

export default SelectTableau;
